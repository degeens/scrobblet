package scrobbler

import (
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/degeens/scrobblet/internal/common"
	"github.com/degeens/scrobblet/internal/events"
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
	bus              *events.Bus

	mu          sync.RWMutex
	lastPlayback *sources.PlaybackState
}

func NewTracker(source sources.Source, playingTrackChan chan common.Track, playedTrackChan chan<- common.TrackedTrack, bus *events.Bus) *Tracker {
	return &Tracker{
		source:           source,
		playingTrackChan: playingTrackChan,
		playedTrackChan:  playedTrackChan,
		bus:              bus,
	}
}

func (t *Tracker) LastPlaybackState() *sources.PlaybackState {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.lastPlayback
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

		t.mu.Lock()
		prevPlayback := t.lastPlayback
		t.lastPlayback = playbackState
		t.mu.Unlock()

		trackEnded := false
		if t.isTrackChange(playbackState, trackedTrack) || t.isTrackReplay(playbackState, trackedTrack) {
			t.sendPlayedTrack(trackedTrack)
			trackedTrack = nil
			trackEnded = true
		}

		if playbackState == nil {
			slog.Debug("No track playing")

			if trackEnded {
				t.publishTrackChange(nil)
			}

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
			t.publishTrackChange(playbackState)

			continue
		}

		if prevPlayback != nil && prevPlayback.IsPlaying != playbackState.IsPlaying {
			t.publishTrackChange(playbackState)
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

func (t *Tracker) publishTrackChange(state *sources.PlaybackState) {
	e := events.Event{Type: events.EventTrackChange}
	if state != nil {
		e.Artist = strings.Join(state.Track.Artists, ", ")
		e.Title = state.Track.Title
		e.Album = state.Track.Album
		e.DurationMs = int64(state.Track.Duration / time.Millisecond)
		e.IsPlaying = state.IsPlaying
		e.PositionMs = int64(state.Position / time.Millisecond)
	}
	t.bus.Publish(e)
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
