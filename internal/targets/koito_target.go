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
	req := t.toPlayingNowSubmitListens(track)

	return t.client.SubmitListens(req)
}

func (t *KoitoTarget) SubmitPlayedTrack(trackedTrack *common.TrackedTrack) error {
	req := t.toSingleSubmitListens(trackedTrack)

	return t.client.SubmitListens(req)
}

func (t *KoitoTarget) toPlayingNowSubmitListens(track *common.Track) *koito.SubmitListens {
	artistName := strings.Join(track.Artists, ", ")

	return &koito.SubmitListens{
		ListenType: "playing_now",
		Payload: []koito.Payload{
			{
				TrackMetadata: koito.TrackMetadata{
					ArtistName:  artistName,
					TrackName:   track.Title,
					ReleaseName: track.Album,
					AdditionalInfo: koito.AdditionalInfo{
						ArtistNames: track.Artists,
					},
				},
			},
		},
	}
}

func (t *KoitoTarget) toSingleSubmitListens(trackedTrack *common.TrackedTrack) *koito.SubmitListens {
	listenedAt := trackedTrack.LastUpdatedAt.Unix()
	artistName := strings.Join(trackedTrack.Track.Artists, ", ")

	return &koito.SubmitListens{
		ListenType: "single",
		Payload: []koito.Payload{
			{
				ListenedAt: &listenedAt,
				TrackMetadata: koito.TrackMetadata{
					ArtistName:  artistName,
					TrackName:   trackedTrack.Track.Title,
					ReleaseName: trackedTrack.Track.Album,
					AdditionalInfo: koito.AdditionalInfo{
						ArtistNames: trackedTrack.Track.Artists,
					},
				},
			},
		},
	}
}
