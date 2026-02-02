package scrobbler

import (
	"github.com/degeens/scrobblet/internal/common"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

type Scrobbler struct {
	source  sources.Source
	targets []targets.Target
}

func NewScrobbler(source sources.Source, targets []targets.Target) *Scrobbler {
	return &Scrobbler{
		source:  source,
		targets: targets,
	}
}

func (s *Scrobbler) Start() {
	playingTrackChan := make(chan common.Track, 1)
	playedTrackChan := make(chan common.TrackedTrack, 10)

	tracker := NewTracker(s.source, playingTrackChan, playedTrackChan)
	go tracker.Start()

	submitter := NewSubmitter(s.targets, playingTrackChan, playedTrackChan)
	go submitter.Start()
}
