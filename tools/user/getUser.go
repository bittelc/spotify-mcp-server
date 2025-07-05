package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"spotify-mcp-server/config"

	"github.com/localrivet/gomcp/server"
)

type GetUserArgs struct {
	// No arguments needed for getting current user profile
}

type UserProfile struct {
	ID              string                  `json:"id"`
	DisplayName     string                  `json:"display_name"`
	Email           string                  `json:"email,omitempty"`
	Country         string                  `json:"country,omitempty"`
	Product         string                  `json:"product,omitempty"`
	URI             string                  `json:"uri"`
	ExternalURL     string                  `json:"external_url"`
	Followers       int                     `json:"followers"`
	Images          []UserImage             `json:"images,omitempty"`
	ExplicitContent ExplicitContentSettings `json:"explicit_content,omitempty"`
}

type UserImage struct {
	URL    string `json:"url"`
	Height *int   `json:"height"`
	Width  *int   `json:"width"`
}

type ExplicitContentSettings struct {
	FilterEnabled bool `json:"filter_enabled"`
	FilterLocked  bool `json:"filter_locked"`
}

type SpotifyUserResponse struct {
	ID           string `json:"id"`
	DisplayName  string `json:"display_name"`
	Email        string `json:"email"`
	Country      string `json:"country"`
	Product      string `json:"product"`
	URI          string `json:"uri"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"followers"`
	Images []struct {
		URL    string `json:"url"`
		Height *int   `json:"height"`
		Width  *int   `json:"width"`
	} `json:"images"`
	ExplicitContent struct {
		FilterEnabled bool `json:"filter_enabled"`
		FilterLocked  bool `json:"filter_locked"`
	} `json:"explicit_content"`
	Href string `json:"href"`
	Type string `json:"type"`
}

func HandleGetUser(ctx *server.Context, args GetUserArgs) (UserProfile, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return UserProfile{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create authenticated client
	auth := config.NewSpotifyAuth(cfg)
	client, err := auth.CreateAuthenticatedClient()
	if err != nil {
		return UserProfile{}, fmt.Errorf("failed to create authenticated client: %w", err)
	}

	// Make the request
	url := "https://api.spotify.com/v1/me"
	resp, err := client.Get(url)
	if err != nil {
		return UserProfile{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserProfile{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return UserProfile{}, fmt.Errorf("spotify API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var spotifyResp SpotifyUserResponse
	if err := json.Unmarshal(body, &spotifyResp); err != nil {
		return UserProfile{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Convert images to our format
	var images []UserImage
	for _, img := range spotifyResp.Images {
		images = append(images, UserImage{
			URL:    img.URL,
			Height: img.Height,
			Width:  img.Width,
		})
	}

	// Convert to our return format
	result := UserProfile{
		ID:          spotifyResp.ID,
		DisplayName: spotifyResp.DisplayName,
		Email:       spotifyResp.Email,
		Country:     spotifyResp.Country,
		Product:     spotifyResp.Product,
		URI:         spotifyResp.URI,
		ExternalURL: spotifyResp.ExternalURLs.Spotify,
		Followers:   spotifyResp.Followers.Total,
		Images:      images,
		ExplicitContent: ExplicitContentSettings{
			FilterEnabled: spotifyResp.ExplicitContent.FilterEnabled,
			FilterLocked:  spotifyResp.ExplicitContent.FilterLocked,
		},
	}

	return result, nil
}
