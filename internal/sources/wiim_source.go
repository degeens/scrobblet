package sources

import (
	"time"

	"github.com/degeens/scrobblet/internal/clients/wiim"
	"github.com/degeens/scrobblet/internal/common"
)

type WiiMSource struct {
	healthy         bool
	lastHealthCheck time.Time
	client          *wiim.Client
}

func NewWiiMSource(client *wiim.Client) *WiiMSource {
	return &WiiMSource{
		healthy:         true,
		lastHealthCheck: time.Now().UTC(),
		client:          client,
	}
}

func (s *WiiMSource) Healthy() (bool, time.Time) {
	return s.healthy, s.lastHealthCheck
}

func (s *WiiMSource) SourceType() SourceType {
	return SourceWiiM
}

func (s *WiiMSource) GetPlaybackState() (*PlaybackState, error) {
	playerStatus, err := s.client.GetPlayerStatus()
	if err != nil {
		s.healthy = false
		s.lastHealthCheck = time.Now().UTC()
		return nil, err
	}

	playbackState := wiimToPlaybackState(playerStatus)

	s.healthy = true
	s.lastHealthCheck = time.Now().UTC()
	return playbackState, nil
}

func wiimToPlaybackState(playerStatus *wiim.PlayerStatus) *PlaybackState {
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
