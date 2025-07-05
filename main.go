package main

import (
	"log"
	"log/slog"
	"os"

	"spotify-mcp-server/tools/playlist"
	"spotify-mcp-server/tools/user"

	"github.com/joho/godotenv"
	"github.com/localrivet/gomcp/server"
)

func main() {
	// Load .env file if it exists (for development)
	if err := godotenv.Load(); err != nil {
		// Don't fail if .env doesn't exist - it's optional
		log.Printf("No .env file found or error loading it: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	s := server.NewServer("spotify_mcp_server",
		server.WithLogger(logger),
	).AsStdio()

	s.Tool("get_user", "get own user's profile and data",
		user.HandleGetUser)
	s.Tool("create_playlist", "create a new playlist",
		playlist.HandleCreatePlaylist)

	if err := s.Run(); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
