package targets

import (
	"strings"

	"github.com/degeens/scrobblet/internal/clients/listenbrainz"
	"github.com/degeens/scrobblet/internal/common"
)

type ListenBrainzTarget struct {
	client *listenbrainz.Client
}

func NewListenBrainzTarget(client *listenbrainz.Client) *ListenBrainzTarget {
	return &ListenBrainzTarget{
		client: client,
	}
}

func (t *ListenBrainzTarget) SubmitPlayingTrack(track *common.Track) error {
	req := t.toPlayingNowSubmitListens(track)

	return t.client.SubmitListens(req)
}

func (t *ListenBrainzTarget) SubmitPlayedTrack(trackedTrack *common.TrackedTrack) error {
	req := t.toSingleSubmitListens(trackedTrack)

	return t.client.SubmitListens(req)
}

func (t *ListenBrainzTarget) toPlayingNowSubmitListens(track *common.Track) *listenbrainz.SubmitListens {
	artistName := strings.Join(track.Artists, ", ")

	return &listenbrainz.SubmitListens{
		ListenType: "playing_now",
		Payload: []listenbrainz.Payload{
			{
				TrackMetadata: listenbrainz.TrackMetadata{
					ArtistName:  artistName,
					TrackName:   track.Title,
					ReleaseName: track.Album,
				},
			},
		},
	}
}

func (t *ListenBrainzTarget) toSingleSubmitListens(trackedTrack *common.TrackedTrack) *listenbrainz.SubmitListens {
	listenedAt := trackedTrack.LastUpdatedAt.Unix()
	artistName := strings.Join(trackedTrack.Track.Artists, ", ")

	return &listenbrainz.SubmitListens{
		ListenType: "single",
		Payload: []listenbrainz.Payload{
			{
				ListenedAt: &listenedAt,
				TrackMetadata: listenbrainz.TrackMetadata{
					ArtistName:  artistName,
					TrackName:   trackedTrack.Track.Title,
					ReleaseName: trackedTrack.Track.Album,
				},
			},
		},
	}
}
