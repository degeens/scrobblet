package scrobbler

import (
	"log/slog"

	"github.com/degeens/scrobblet/internal/common"
	"github.com/degeens/scrobblet/internal/targets"
)

type Submitter struct {
	target           targets.Target
	playingTrackChan <-chan common.Track
	playedTrackChan  <-chan common.TrackedTrack
}

func NewSubmitter(target targets.Target, playingTrackChan <-chan common.Track, playedTrackChan <-chan common.TrackedTrack) *Submitter {
	return &Submitter{
		target:           target,
		playingTrackChan: playingTrackChan,
		playedTrackChan:  playedTrackChan,
	}
}

func (s *Submitter) Start() {
	for {
		select {
		case track := <-s.playingTrackChan:
			slog.Info("Submitting now playing track", track.SlogArgs()...)

			err := s.target.SubmitPlayingTrack(&track)
			if err != nil {
				slog.Error("Failed to submit now playing track", append(track.SlogArgs(), "error", err.Error())...)
				continue
				// todo: retry (with exponential backoff)
			}

			slog.Info("Now playing track submitted", track.SlogArgs()...)
		case trackedTrack := <-s.playedTrackChan:
			if !ShouldScrobble(trackedTrack.Duration, trackedTrack.Track.Duration) {
				slog.Info("Track did not meet scrobble rules, skipping track", trackedTrack.Track.SlogArgs()...)
				continue
			}

			slog.Info("Track met scrobble rules, submitting track", trackedTrack.Track.SlogArgs()...)

			err := s.target.SubmitPlayedTrack(&trackedTrack)
			if err != nil {
				slog.Error("Failed to submit track", append(trackedTrack.Track.SlogArgs(), "error", err.Error())...)
				continue
				// todo: retry (with exponential backoff)
			}

			slog.Info("Track submitted", trackedTrack.Track.SlogArgs()...)
		}
	}
}
