package config

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// TokenResponse represents the response from Spotify's token endpoint
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// SpotifyAuth handles authentication with the Spotify API
type SpotifyAuth struct {
	config      *SpotifyConfig
	accessToken string
	expiresAt   time.Time
}

// NewSpotifyAuth creates a new SpotifyAuth instance
func NewSpotifyAuth(config *SpotifyConfig) *SpotifyAuth {
	return &SpotifyAuth{
		config: config,
	}
}

// GetAccessToken returns a valid access token, refreshing it if necessary
func (s *SpotifyAuth) GetAccessToken() (string, error) {
	// Check if we have a valid token that hasn't expired
	if s.accessToken != "" && time.Now().Before(s.expiresAt) {
		return s.accessToken, nil
	}

	// Token is expired or doesn't exist, get a new one
	token, err := s.requestAccessToken()
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	s.accessToken = token.AccessToken
	// Set expiration time with a small buffer to avoid edge cases
	s.expiresAt = time.Now().Add(time.Duration(token.ExpiresIn-60) * time.Second)

	return s.accessToken, nil
}

// requestAccessToken requests a new access token using the client credentials flow
func (s *SpotifyAuth) requestAccessToken() (*TokenResponse, error) {
	// Prepare the request body
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	// Create the authorization header
	auth := s.config.ClientID + ":" + s.config.ClientSecret
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	// Create the HTTP request
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Basic "+encodedAuth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication failed with status code: %d", resp.StatusCode)
	}

	// Parse the response
	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &tokenResponse, nil
}

// IsTokenValid checks if the current token is valid and not expired
func (s *SpotifyAuth) IsTokenValid() bool {
	return s.accessToken != "" && time.Now().Before(s.expiresAt)
}

// ClearToken clears the current token, forcing a refresh on next access
func (s *SpotifyAuth) ClearToken() {
	s.accessToken = ""
	s.expiresAt = time.Time{}
}

// CreateAuthenticatedClient creates an HTTP client with authentication headers
func (s *SpotifyAuth) CreateAuthenticatedClient() (*http.Client, error) {
	token, err := s.GetAccessToken()
	if err != nil {
		return nil, err
	}

	// Create a custom transport that adds the Authorization header
	transport := &AuthTransport{
		Base:        http.DefaultTransport,
		AccessToken: token,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}, nil
}

// AuthTransport is a custom HTTP transport that adds authorization headers
type AuthTransport struct {
	Base        http.RoundTripper
	AccessToken string
}

// RoundTrip implements the RoundTripper interface
func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	reqCopy := req.Clone(req.Context())
	reqCopy.Header.Set("Authorization", "Bearer "+t.AccessToken)

	// Use the base transport to make the request
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}

	return base.RoundTrip(reqCopy)
}
