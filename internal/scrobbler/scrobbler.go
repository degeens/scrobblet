package scrobbler

import (
	"github.com/degeens/scrobblet/internal/common"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

type Scrobbler struct {
	source sources.Source
	target targets.Target
}

func NewScrobbler(source sources.Source, target targets.Target) *Scrobbler {
	return &Scrobbler{
		source: source,
		target: target,
	}
}

func (s *Scrobbler) Start() {
	playingTrackChan := make(chan common.Track, 1)
	playedTrackChan := make(chan common.TrackedTrack, 10)

	tracker := NewTracker(s.source, playingTrackChan, playedTrackChan)
	go tracker.Start()

	submitter := NewSubmitter(s.target, playingTrackChan, playedTrackChan)
	go submitter.Start()
}
