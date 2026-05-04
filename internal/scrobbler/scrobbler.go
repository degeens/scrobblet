package scrobbler

import (
	"github.com/degeens/scrobblet/internal/common"
	"github.com/degeens/scrobblet/internal/events"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

type Scrobbler struct {
	targets          []targets.Target
	bus              *events.Bus
	tracker          *Tracker
	playingTrackChan chan common.Track
	playedTrackChan  chan common.TrackedTrack
}

func NewScrobbler(source sources.Source, targets []targets.Target, bus *events.Bus) *Scrobbler {
	playingTrackChan := make(chan common.Track, 1)
	playedTrackChan := make(chan common.TrackedTrack, 10)

	return &Scrobbler{
		targets:          targets,
		bus:              bus,
		tracker:          NewTracker(source, playingTrackChan, playedTrackChan, bus),
		playingTrackChan: playingTrackChan,
		playedTrackChan:  playedTrackChan,
	}
}

func (s *Scrobbler) Start() {
	go s.tracker.Start()

	submitter := NewSubmitter(s.targets, s.playingTrackChan, s.playedTrackChan, s.bus)
	go submitter.Start()
}

func (s *Scrobbler) LastPlaybackState() *sources.PlaybackState {
	return s.tracker.LastPlaybackState()
}
