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
	Healthy() (bool, time.Time)
	SourceType() SourceType
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
		client, err := spotify.NewClient(clientsConfig.Spotify.ClientID, clientsConfig.Spotify.ClientSecret, clientsConfig.Spotify.RedirectURL, clientsConfig.Spotify.DataPath)
		if err != nil {
			return nil, nil, err
		}
		return client, NewSpotifySource(client), nil
	default:
		return nil, nil, fmt.Errorf("unknown source type: %s", sourceType)
	}
}
