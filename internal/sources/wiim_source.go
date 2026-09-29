package sources

import (
	"strings"
	"time"

	"github.com/degeens/scrobblet/internal/clients/wiim"
	"github.com/degeens/scrobblet/internal/common"
)

type WiiMSource struct {
	health *common.HealthStatus
	client *wiim.Client
}

func NewWiiMSource(client *wiim.Client) *WiiMSource {
	return &WiiMSource{
		health: common.NewHealthStatus(),
		client: client,
	}
}

func (s *WiiMSource) Healthy() (bool, time.Time) {
	return s.health.Get()
}

func (s *WiiMSource) SourceType() SourceType {
	return SourceWiiM
}

func (s *WiiMSource) GetPlaybackState() (*PlaybackState, error) {
	playerStatus, err := s.client.GetPlayerStatus()
	if err != nil {
		s.health.Set(false)
		return nil, err
	}

	playbackState := wiimToPlaybackState(playerStatus)

	s.health.Set(true)
	return playbackState, nil
}

func wiimToPlaybackState(playerStatus *wiim.PlayerStatus) *PlaybackState {
	const unknown = "Unknown"

	// For example, WiiM reports "Unknown" metadata for optical inputs, such as TV audio, which should not be scrobbled.
	if strings.EqualFold(playerStatus.Title, unknown) ||
		strings.EqualFold(playerStatus.Artist, unknown) ||
		strings.EqualFold(playerStatus.Album, unknown) {
		return nil
	}

	return &PlaybackState{
		Track: &common.Track{
			Artists:  []string{playerStatus.Artist},
			Title:    playerStatus.Title,
			Album:    playerStatus.Album,
			Duration: time.Duration(playerStatus.TotalLength) * time.Millisecond,
		},
		Position:  time.Duration(playerStatus.CurrentPosition) * time.Millisecond,
		Timestamp: time.Now().UTC(),
	}
}
