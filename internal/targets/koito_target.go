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

func (t *KoitoTarget) SubmitPlayingTrack(track *common.Track) error {
	req := toPlayingNowSubmitListens(track)

	return t.client.SubmitListens(req)
}

func (t *KoitoTarget) SubmitPlayedTrack(trackedTrack *common.TrackedTrack) error {
	req := toSingleSubmitListens(trackedTrack)

	return t.client.SubmitListens(req)
}

func toPlayingNowSubmitListens(track *common.Track) *koito.SubmitListens {
	artistName := strings.Join(track.Artists, ", ")

	return &koito.SubmitListens{
		ListenType: "playing_now",
		Payload: []koito.Payload{
			{
				TrackMetadata: koito.TrackMetadata{
					ArtistName:  artistName,
					TrackName:   track.Title,
					ReleaseName: track.Album,
				},
			},
		},
	}
}

func toSingleSubmitListens(trackedTrack *common.TrackedTrack) *koito.SubmitListens {
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
