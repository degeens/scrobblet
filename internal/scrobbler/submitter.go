package scrobbler

import (
	"log/slog"

	"github.com/degeens/scrobblet/internal/common"
	"github.com/degeens/scrobblet/internal/events"
	"github.com/degeens/scrobblet/internal/targets"
)

type Submitter struct {
	targets          []targets.Target
	playingTrackChan <-chan common.Track
	playedTrackChan  <-chan common.TrackedTrack
	bus              *events.Bus
}

func NewSubmitter(targets []targets.Target, playingTrackChan <-chan common.Track, playedTrackChan <-chan common.TrackedTrack, bus *events.Bus) *Submitter {
	return &Submitter{
		targets:          targets,
		playingTrackChan: playingTrackChan,
		playedTrackChan:  playedTrackChan,
		bus:              bus,
	}
}

func (s *Submitter) Start() {
	for {
		select {
		case track := <-s.playingTrackChan:
			slog.Info("Submitting now playing track", track.SlogArgs()...)

			for _, target := range s.targets {
				err := target.SubmitPlayingTrack(&track)
				if err != nil {
					slog.Error("Failed to submit now playing track", append(track.SlogArgs(), "target", target.TargetType(), "error", err.Error())...)
					s.bus.Publish(events.Event{Type: events.EventNowPlaying, Target: string(target.TargetType()), Success: false})
					continue
					// todo: retry (with exponential backoff)
				}

				slog.Info("Now playing track submitted", append(track.SlogArgs(), "target", target.TargetType())...)
				s.bus.Publish(events.Event{Type: events.EventNowPlaying, Target: string(target.TargetType()), Success: true})
			}
		case trackedTrack := <-s.playedTrackChan:
			if !ShouldScrobble(trackedTrack.Duration, trackedTrack.Track.Duration) {
				slog.Info("Track did not meet scrobble rules, skipping track", trackedTrack.Track.SlogArgs()...)
				continue
			}

			slog.Info("Track met scrobble rules, submitting track", trackedTrack.Track.SlogArgs()...)

			for _, target := range s.targets {
				err := target.SubmitPlayedTrack(&trackedTrack)
				if err != nil {
					slog.Error("Failed to submit track", append(trackedTrack.Track.SlogArgs(), "target", target.TargetType(), "error", err.Error())...)
					s.bus.Publish(events.Event{Type: events.EventScrobble, Target: string(target.TargetType()), Success: false})
					continue
					// todo: retry (with exponential backoff)
				}

				slog.Info("Track submitted", append(trackedTrack.Track.SlogArgs(), "target", target.TargetType())...)
				s.bus.Publish(events.Event{Type: events.EventScrobble, Target: string(target.TargetType()), Success: true})
			}
		}
	}
}
