package sources

import (
	"time"

	"github.com/degeens/scrobblet/internal/clients/spotify"
	"github.com/degeens/scrobblet/internal/common"
)

type SpotifySource struct {
	health *common.HealthStatus
	client *spotify.Client
}

func NewSpotifySource(client *spotify.Client) *SpotifySource {
	return &SpotifySource{
		health: common.NewHealthStatus(),
		client: client,
	}
}

func (s *SpotifySource) Healthy() (bool, time.Time) {
	return s.health.Get()
}

func (s *SpotifySource) SourceType() SourceType {
	return SourceSpotify
}

func (s *SpotifySource) GetPlaybackState() (*PlaybackState, error) {
	currentlyPlaying, err := s.client.GetCurrentlyPlayingTrack()
	if err != nil {
		s.health.Set(false)
		return nil, err
	}

	playbackState := spotifyToPlaybackState(currentlyPlaying)

	s.health.Set(true)
	return playbackState, nil
}

func spotifyToPlaybackState(currentlyPlaying *spotify.CurrentlyPlayingTrack) *PlaybackState {
	if currentlyPlaying == nil {
		return nil
	}

	artists := make([]string, len(currentlyPlaying.Item.Artists))
	for i, artist := range currentlyPlaying.Item.Artists {
		artists[i] = artist.Name
	}

	return &PlaybackState{
		Track: &common.Track{
			Artists:     artists,
			Title:       currentlyPlaying.Item.Name,
			Album:       currentlyPlaying.Item.Album.Name,
			TrackNumber: currentlyPlaying.Item.TrackNumber,
			Duration:    time.Duration(currentlyPlaying.Item.Duration) * time.Millisecond,
		},
		Position:  time.Duration(currentlyPlaying.Progress) * time.Millisecond,
		Timestamp: time.Now().UTC(),
	}
}
