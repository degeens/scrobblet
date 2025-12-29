package scrobbler

import (
	"log/slog"
	"time"

	"github.com/degeens/scrobblet/internal/common"
	"github.com/degeens/scrobblet/internal/sources"
)

const (
	pollInterval   = 10 * time.Second
	driftTolerance = 2 * time.Second
)

type Tracker struct {
	source    sources.Source
	trackChan chan<- common.TrackedTrack
}

func NewTracker(source sources.Source, trackChan chan<- common.TrackedTrack) *Tracker {
	return &Tracker{
		source:    source,
		trackChan: trackChan,
	}
}

func (t *Tracker) Start() {
	ticker := time.NewTicker(pollInterval)

	var trackedTrack *common.TrackedTrack

	for range ticker.C {
		playbackState, err := t.source.GetPlaybackState()
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		// Handle previous track
		if trackedTrack != nil && (playbackState == nil || !playbackState.Track.Equals(trackedTrack.Track)) {
			select {
			case t.trackChan <- *trackedTrack:
				slog.Info("Track added to queue", trackedTrack.Track.SlogArgs()...)
			default:
				slog.Warn("Track queue is full, skipping track", trackedTrack.Track.SlogArgs()...)
			}

			trackedTrack = nil
		}

		// Handle current track
		if playbackState == nil {
			slog.Debug("No track playing")
			continue
		}

		now := time.Now().UTC()

		// Start tracking new track
		if trackedTrack == nil {
			trackedTrack = &common.TrackedTrack{
				Track:         playbackState.Track,
				LastPosition:  playbackState.Position,
				Duration:      0,
				StartedAt:     now,
				LastUpdatedAt: now,
			}

			slog.Info("Track is being tracked", trackedTrack.SlogArgs()...)

			continue
		}

		// Continue tracking existing track
		positionDiff := playbackState.Position - trackedTrack.LastPosition
		timeDiff := playbackState.Timestamp.Sub(trackedTrack.LastUpdatedAt)

		drift := positionDiff - timeDiff
		if drift < 0 {
			drift = -drift
		}

		if drift > driftTolerance {
			trackedTrack.LastPosition = playbackState.Position

			slog.Info("Seek or pause detected", trackedTrack.SlogArgs()...)
		} else {
			trackedTrack.Duration += positionDiff
			trackedTrack.LastPosition = playbackState.Position
			trackedTrack.LastUpdatedAt = now

			slog.Info("Track is being tracked", trackedTrack.SlogArgs()...)
		}
	}
}
