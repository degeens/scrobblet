package sources

import (
	"fmt"
	"time"

	"github.com/degeens/scrobblet/internal/clients/spotify"
	"github.com/degeens/scrobblet/internal/clients/wiim"
	"github.com/degeens/scrobblet/internal/common"
)

type SourceType string

const (
	SourceSpotify SourceType = "Spotify"
	SourceWiiM    SourceType = "WiiM"
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

// Spec describes how to construct a source. The config field matching Type
// must be set.
type Spec struct {
	Type    SourceType
	Spotify *spotify.Config
	WiiM    *wiim.Config
}

func New(spec Spec) (any, Source, error) {
	switch spec.Type {
	case SourceSpotify:
		if spec.Spotify == nil {
			return nil, nil, fmt.Errorf("missing Spotify configuration for source type: %s", spec.Type)
		}
		client, err := spotify.NewClient(spec.Spotify.ClientID, spec.Spotify.ClientSecret, spec.Spotify.RedirectURL, spec.Spotify.DataPath)
		if err != nil {
			return nil, nil, err
		}
		return client, NewSpotifySource(client), nil
	case SourceWiiM:
		if spec.WiiM == nil {
			return nil, nil, fmt.Errorf("missing WiiM configuration for source type: %s", spec.Type)
		}
		client := wiim.NewClient(spec.WiiM.BaseURL)
		return client, NewWiiMSource(client), nil
	default:
		return nil, nil, fmt.Errorf("unknown source type: %s", spec.Type)
	}
}
