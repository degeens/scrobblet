package scrobbler

import (
	"log/slog"
	"time"

	"github.com/degeens/scrobblet/internal/common"
	"github.com/degeens/scrobblet/internal/sources"
)

const (
	activePollInterval   = 10 * time.Second
	inactivePollInterval = 30 * time.Second
	inactivityThreshold  = 5 * time.Minute
	driftTolerance       = 2 * time.Second
)

type Tracker struct {
	source           sources.Source
	playingTrackChan chan common.Track
	playedTrackChan  chan<- common.TrackedTrack
}

func NewTracker(source sources.Source, playingTrackChan chan common.Track, playedTrackChan chan<- common.TrackedTrack) *Tracker {
	return &Tracker{
		source:           source,
		playingTrackChan: playingTrackChan,
		playedTrackChan:  playedTrackChan,
	}
}

func (t *Tracker) Start() {
	lastActivity := time.Now().UTC()
	pollInterval := activePollInterval
	ticker := time.NewTicker(pollInterval)

	var trackedTrack *common.TrackedTrack

	for range ticker.C {
		playbackState, err := t.source.GetPlaybackState()
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		if t.isTrackChange(playbackState, trackedTrack) || t.isTrackReplay(playbackState, trackedTrack) {
			t.sendPlayedTrack(trackedTrack)
			trackedTrack = nil
		}

		if playbackState == nil {
			slog.Debug("No track playing")

			pollInterval = t.switchToInactivePollingIntervalIfNeeded(ticker, pollInterval, lastActivity)

			continue
		}

		now := time.Now().UTC()

		lastActivity = now
		pollInterval = t.switchToActivePollingIntervalIfNeeded(ticker, pollInterval)

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

			t.sendPlayingTrack(trackedTrack.Track)

			continue
		}

		positionDiff := playbackState.Position - trackedTrack.LastPosition
		timeDiff := playbackState.Timestamp.Sub(trackedTrack.LastUpdatedAt)

		// Continue tracking existing track
		if t.isNormalPlayback(positionDiff, timeDiff) {
			trackedTrack.Duration += positionDiff
			trackedTrack.LastPosition = playbackState.Position
			trackedTrack.LastUpdatedAt = now

			slog.Info("Track is being tracked", trackedTrack.SlogArgs()...)
		} else {
			trackedTrack.LastPosition = playbackState.Position
			trackedTrack.LastUpdatedAt = now

			slog.Info("Seek or pause detected", trackedTrack.SlogArgs()...)
		}
	}
}

func (t *Tracker) isTrackChange(playbackState *sources.PlaybackState, trackedTrack *common.TrackedTrack) bool {
	return trackedTrack != nil &&
		(playbackState == nil || !playbackState.Track.Equals(trackedTrack.Track))
}

func (t *Tracker) isTrackReplay(playbackState *sources.PlaybackState, trackedTrack *common.TrackedTrack) bool {
	return trackedTrack != nil &&
		playbackState != nil &&
		playbackState.Track.Equals(trackedTrack.Track) &&
		playbackState.Position <= activePollInterval &&
		trackedTrack.LastPosition > activePollInterval
}

func (t *Tracker) isNormalPlayback(positionDiff, timeDiff time.Duration) bool {
	drift := positionDiff - timeDiff
	if drift < 0 {
		drift = -drift
	}

	return drift < driftTolerance
}

func (t *Tracker) switchToInactivePollingIntervalIfNeeded(ticker *time.Ticker, currentInterval time.Duration, lastActivityTime time.Time) time.Duration {
	if time.Since(lastActivityTime) > inactivityThreshold && currentInterval != inactivePollInterval {
		ticker.Reset(inactivePollInterval)

		slog.Info("Switched to inactive polling interval", "interval", inactivePollInterval/time.Second)

		return inactivePollInterval
	}

	return currentInterval
}

func (t *Tracker) switchToActivePollingIntervalIfNeeded(ticker *time.Ticker, currentInterval time.Duration) time.Duration {
	if currentInterval != activePollInterval {
		ticker.Reset(activePollInterval)

		slog.Info("Switched to active polling interval", "interval", activePollInterval/time.Second)

		return activePollInterval
	}

	return currentInterval
}

func (t *Tracker) sendPlayedTrack(trackedTrack *common.TrackedTrack) {
	select {
	case t.playedTrackChan <- *trackedTrack:
		slog.Info("Track added to queue", trackedTrack.Track.SlogArgs()...)
	default:
		slog.Warn("Track queue is full, skipping track", trackedTrack.Track.SlogArgs()...)
	}
}

func (t *Tracker) sendPlayingTrack(track *common.Track) {
	// Drain any pending track to ensure the latest track wins
	select {
	case <-t.playingTrackChan:
	default:
	}

	select {
	case t.playingTrackChan <- *track:
	default:
		// This should never happen after draining with buffer size 1
		slog.Error("Failed to send now playing track to channel", track.SlogArgs()...)
	}
}
