package targets

import (
	"time"

	"github.com/degeens/scrobblet/internal/clients/csv"
	"github.com/degeens/scrobblet/internal/common"
)

type CSVTarget struct {
	healthy         bool
	lastHealthCheck time.Time
	client          *csv.Client
}

func NewCSVTarget(client *csv.Client) *CSVTarget {
	return &CSVTarget{
		healthy:         true,
		lastHealthCheck: time.Now().UTC(),
		client:          client,
	}
}

func (t *CSVTarget) Healthy() (bool, time.Time) {
	return t.healthy, t.lastHealthCheck
}

func (t *CSVTarget) TargetType() TargetType {
	return TargetCSV
}

func (t *CSVTarget) SubmitPlayingTrack(track *common.Track) error {
	// Only submit completed scrobbles
	return nil
}

func (t *CSVTarget) SubmitPlayedTrack(trackedTrack *common.TrackedTrack) error {
	err := t.client.WriteScrobble(trackedTrack)
	if err != nil {
		t.healthy = false
		t.lastHealthCheck = time.Now().UTC()
		return err
	}

	t.healthy = true
	t.lastHealthCheck = time.Now().UTC()
	return nil
}
