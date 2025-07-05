package playlist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"spotify-mcp-server/config"

	"github.com/localrivet/gomcp/server"
)

type CreatePlaylistArgs struct {
	UserID string `json:"user_id"`
	// Name          string `json:"name"`
	// Description   string `json:"description,omitempty"`
	// Public        *bool  `json:"public,omitempty"`
	// Collaborative *bool  `json:"collaborative,omitempty"`
}

type PlaylistMetadata struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Public      bool   `json:"public"`
	URI         string `json:"uri"`
	ExternalURL string `json:"external_url"`
}

type SpotifyPlaylistRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Public        *bool  `json:"public,omitempty"`
	Collaborative *bool  `json:"collaborative,omitempty"`
}

type SpotifyPlaylistResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Public       bool   `json:"public"`
	URI          string `json:"uri"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

func HandleCreatePlaylist(ctx *server.Context, args CreatePlaylistArgs) (PlaylistMetadata, error) {
	// Validate required fields
	if args.UserID == "" {
		return PlaylistMetadata{}, fmt.Errorf("user_id is required")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return PlaylistMetadata{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Get access token from config
	apiKey := cfg.ClientSecret

	// Prepare the request body
	requestBody := SpotifyPlaylistRequest{
		Name: time.Now().String(),
		// Description:   args.Description,
		// Public:        args.Public,
		// Collaborative: args.Collaborative,
	}

	// If public is not explicitly set, default to true
	if requestBody.Public == nil {
		defaultPublic := true
		requestBody.Public = &defaultPublic
	}

	// If collaborative is not explicitly set, default to false
	if requestBody.Collaborative == nil {
		defaultCollaborative := false
		requestBody.Collaborative = &defaultCollaborative
	}

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return PlaylistMetadata{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create the HTTP request
	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", args.UserID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return PlaylistMetadata{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return PlaylistMetadata{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PlaylistMetadata{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusCreated {
		return PlaylistMetadata{}, fmt.Errorf("spotify API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var spotifyResp SpotifyPlaylistResponse
	if err := json.Unmarshal(body, &spotifyResp); err != nil {
		return PlaylistMetadata{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Convert to our return format
	result := PlaylistMetadata{
		ID:          spotifyResp.ID,
		Name:        spotifyResp.Name,
		Description: spotifyResp.Description,
		Public:      spotifyResp.Public,
		URI:         spotifyResp.URI,
		ExternalURL: spotifyResp.ExternalURLs.Spotify,
	}

	return result, nil
}
