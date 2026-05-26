package scrobbler

import (
	"github.com/degeens/scrobblet/internal/common"
	"github.com/degeens/scrobblet/internal/metrics"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

type Scrobbler struct {
	source  sources.Source
	targets []targets.Target
	metrics *metrics.Metrics
}

func NewScrobbler(source sources.Source, targets []targets.Target, metrics *metrics.Metrics) *Scrobbler {
	return &Scrobbler{
		source:  source,
		targets: targets,
		metrics: metrics,
	}
}

func (s *Scrobbler) Start() {
	playingTrackChan := make(chan common.Track, 1)
	playedTrackChan := make(chan common.TrackedTrack, 10)

	tracker := NewTracker(s.source, playingTrackChan, playedTrackChan, s.metrics)
	go tracker.Start()

	submitter := NewSubmitter(s.targets, playingTrackChan, playedTrackChan, s.metrics)
	go submitter.Start()
}
