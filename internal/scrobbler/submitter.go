package scrobbler

import (
	"log/slog"

	"github.com/degeens/scrobblet/internal/common"
	"github.com/degeens/scrobblet/internal/targets"
)

type Submitter struct {
	target    targets.Target
	trackChan <-chan common.TrackedTrack
}

func NewSubmitter(target targets.Target, trackChan <-chan common.TrackedTrack) *Submitter {
	return &Submitter{
		target:    target,
		trackChan: trackChan,
	}
}

func (s *Submitter) Start() {
	for trackedTrack := range s.trackChan {
		reachedScrobbleTreshold := common.HasReachedScrobbleThreshold(trackedTrack.Duration, trackedTrack.Track.Duration)
		if !reachedScrobbleTreshold {
			slog.Info("Track did not reach scrobble threshold, skipping track", trackedTrack.Track.SlogArgs()...)
			continue
		}

		slog.Info("Track reached scrobble threshold, submitting track", trackedTrack.Track.SlogArgs()...)

		err := s.target.SubmitTrack(&trackedTrack)
		if err != nil {
			slog.Error("Failed to submit track", append(trackedTrack.Track.SlogArgs(), "error", err)...)
			continue
			// todo: retry (with exponential backoff)
		}

		slog.Info("Track submitted", trackedTrack.Track.SlogArgs()...)
	}
}
