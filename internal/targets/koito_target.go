package targets

import (
	"strings"

	"github.com/degeens/scrobblet/internal/clients/koito"
	"github.com/degeens/scrobblet/internal/common"
)

type KoitoTarget struct {
	client *koito.Client
}

func NewKoitoTarget(client *koito.Client) *KoitoTarget {
	return &KoitoTarget{
		client: client,
	}
}

func (t *KoitoTarget) SubmitTrack(trackedTrack *common.TrackedTrack) error {
	req := toSubmitListens(trackedTrack)

	err := t.client.SubmitListens(req)

	return err
}

func toSubmitListens(trackedTrack *common.TrackedTrack) *koito.SubmitListens {
	artistName := strings.Join(trackedTrack.Track.Artists, ", ")

	return &koito.SubmitListens{
		ListenType: "single",
		Payload: []koito.Payload{
			{
				ListenedAt: trackedTrack.LastUpdatedAt.Unix(),
				TrackMetadata: koito.TrackMetadata{
					ArtistName:  artistName,
					TrackName:   trackedTrack.Track.Title,
					ReleaseName: trackedTrack.Track.Album,
				},
			},
		},
	}
}
