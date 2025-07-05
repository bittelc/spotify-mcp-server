package main

import (
	"log"
	"log/slog"
	"os"

	"spotify-mcp-server/tools/playlist"

	"github.com/localrivet/gomcp/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	s := server.NewServer("spotify_mcp_server",
		server.WithLogger(logger),
	).AsStdio()

	s.Tool("create_playlist", "create a new playlist",
		playlist.HandleCreatePlaylist)

	if err := s.Run(); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
