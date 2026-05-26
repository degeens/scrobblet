package scrobbler

import (
	"log/slog"
	"time"

	"github.com/degeens/scrobblet/internal/common"
	"github.com/degeens/scrobblet/internal/metrics"
	"github.com/degeens/scrobblet/internal/targets"
)

type Submitter struct {
	targets          []targets.Target
	playingTrackChan <-chan common.Track
	playedTrackChan  <-chan common.TrackedTrack
	metrics          *metrics.Metrics
}

func NewSubmitter(targets []targets.Target, playingTrackChan <-chan common.Track, playedTrackChan <-chan common.TrackedTrack, metrics *metrics.Metrics) *Submitter {
	return &Submitter{
		targets:          targets,
		playingTrackChan: playingTrackChan,
		playedTrackChan:  playedTrackChan,
		metrics:          metrics,
	}
}

func (s *Submitter) Start() {
	for {
		select {
		case track := <-s.playingTrackChan:
			slog.Info("Submitting now playing track", track.SlogArgs()...)

			for _, target := range s.targets {
				start := time.Now()
				err := target.SubmitPlayingTrack(&track)
				s.metrics.SubmitDuration.WithLabelValues(string(target.TargetType()), metrics.KindNowPlaying).Observe(time.Since(start).Seconds())

				if err != nil {
					s.metrics.NowPlayingSubmits.WithLabelValues(string(target.TargetType()), metrics.StatusFailure).Inc()
					slog.Error("Failed to submit now playing track", append(track.SlogArgs(), "target", target.TargetType(), "error", err.Error())...)
					continue
					// todo: retry (with exponential backoff)
				}

				s.metrics.NowPlayingSubmits.WithLabelValues(string(target.TargetType()), metrics.StatusSuccess).Inc()
				slog.Info("Now playing track submitted", append(track.SlogArgs(), "target", target.TargetType())...)
			}
		case trackedTrack := <-s.playedTrackChan:
			if !ShouldScrobble(trackedTrack.Duration, trackedTrack.Track.Duration) {
				s.metrics.ScrobbleEvaluations.WithLabelValues(metrics.NotMet).Inc()
				slog.Info("Track did not meet scrobble rules, skipping track", trackedTrack.Track.SlogArgs()...)
				continue
			}

			s.metrics.ScrobbleEvaluations.WithLabelValues(metrics.Met).Inc()
			slog.Info("Track met scrobble rules, submitting track", trackedTrack.Track.SlogArgs()...)

			for _, target := range s.targets {
				start := time.Now()
				err := target.SubmitPlayedTrack(&trackedTrack)
				s.metrics.SubmitDuration.WithLabelValues(string(target.TargetType()), metrics.KindScrobble).Observe(time.Since(start).Seconds())

				if err != nil {
					s.metrics.ScrobbleSubmits.WithLabelValues(string(target.TargetType()), metrics.StatusFailure).Inc()
					slog.Error("Failed to submit track", append(trackedTrack.Track.SlogArgs(), "target", target.TargetType(), "error", err.Error())...)
					continue
					// todo: retry (with exponential backoff)
				}

				s.metrics.ScrobbleSubmits.WithLabelValues(string(target.TargetType()), metrics.StatusSuccess).Inc()
				slog.Info("Track submitted", append(trackedTrack.Track.SlogArgs(), "target", target.TargetType())...)
			}
		}
	}
}
