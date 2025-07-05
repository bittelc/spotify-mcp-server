package playlist

import "github.com/localrivet/gomcp/server"

type CreatePlaylistArgs struct{}

type PlaylistMetadata struct {
	ID   string
	Name string
}

func HandleCreatePlaylist(ctx *server.Context, args CreatePlaylistArgs) (PlaylistMetadata, error) {
	return PlaylistMetadata{
		ID:   "12345",
		Name: "My Playlist",
	}, nil
}
