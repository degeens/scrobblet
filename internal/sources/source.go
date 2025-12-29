package sources

import (
	"time"

	"github.com/degeens/scrobblet/internal/common"
)

type Source interface {
	GetPlaybackState() (*PlaybackState, error)
}

type PlaybackState struct {
	Track     *common.Track
	Position  time.Duration
	Timestamp time.Time
}
