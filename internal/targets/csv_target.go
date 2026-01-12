package targets

import (
	"github.com/degeens/scrobblet/internal/clients/csv"
	"github.com/degeens/scrobblet/internal/common"
)

type CSVTarget struct {
	client *csv.Client
}

func NewCSVTarget(client *csv.Client) *CSVTarget {
	return &CSVTarget{
		client: client,
	}
}

func (t *CSVTarget) SubmitPlayingTrack(track *common.Track) error {
	// Only submit completed scrobbles
	return nil
}

func (t *CSVTarget) SubmitPlayedTrack(trackedTrack *common.TrackedTrack) error {
	return t.client.WriteScrobble(trackedTrack)
}
