package config

import (
	"fmt"
	"os"
)

type SpotifyConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func Load() (*SpotifyConfig, error) {
	cfg := &SpotifyConfig{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("SPOTIFY_REDIRECT_URI"),
	}

	// Validate required fields
	if cfg.ClientSecret == "" {
		return nil, fmt.Errorf("SPOTIFY_CLIENT_SECRET is required")
	}
	if cfg.ClientID == "" {
		return nil, fmt.Errorf("SPOTIFY_CLIENT_ID is required")
	}
	if cfg.RedirectURI == "" {
		return nil, fmt.Errorf("SPOTIFY_REDIRECT_URI is required")
	}

	return cfg, nil
}
