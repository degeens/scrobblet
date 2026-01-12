package sources

import (
	"fmt"
	"time"

	"github.com/degeens/scrobblet/internal/clients"
	"github.com/degeens/scrobblet/internal/clients/spotify"
	"github.com/degeens/scrobblet/internal/common"
)

type SourceType string

const (
	SourceSpotify SourceType = "Spotify"
)

type Source interface {
	GetPlaybackState() (*PlaybackState, error)
}

type PlaybackState struct {
	Track     *common.Track
	Position  time.Duration
	Timestamp time.Time
}

func New(sourceType SourceType, clientsConfig clients.Config) (any, Source, error) {
	switch sourceType {
	case SourceSpotify:
		client := spotify.NewClient(clientsConfig.Spotify.ClientID, clientsConfig.Spotify.ClientSecret, clientsConfig.Spotify.RedirectURL, clientsConfig.Spotify.DataPath)
		return client, NewSpotifySource(client), nil
	default:
		return nil, nil, fmt.Errorf("unknown source type: %s", sourceType)
	}
}
