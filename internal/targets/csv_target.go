package targets

import (
	"time"

	"github.com/degeens/scrobblet/internal/clients/csv"
	"github.com/degeens/scrobblet/internal/common"
)

type CSVTarget struct {
	health *common.HealthStatus
	client *csv.Client
}

func NewCSVTarget(client *csv.Client) *CSVTarget {
	return &CSVTarget{
		health: common.NewHealthStatus(),
		client: client,
	}
}

func (t *CSVTarget) Healthy() (bool, time.Time) {
	return t.health.Get()
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
		t.health.Set(false)
		return err
	}

	t.health.Set(true)
	return nil
}
