package targets

import (
	"strconv"
	"strings"
	"time"

	"github.com/degeens/scrobblet/internal/clients/listenbrainz"
	"github.com/degeens/scrobblet/internal/common"
)

type ListenBrainzTarget struct {
	healthy          bool
	lastHealthCheck  time.Time
	targetType       TargetType
	client           *listenbrainz.Client
	scrobbletVersion string
}

func NewListenBrainzTarget(targetType TargetType, client *listenbrainz.Client, scrobbletVersion string) *ListenBrainzTarget {
	return &ListenBrainzTarget{
		healthy:          true,
		lastHealthCheck:  time.Now().UTC(),
		targetType:       targetType,
		client:           client,
		scrobbletVersion: scrobbletVersion,
	}
}

func (t *ListenBrainzTarget) Healthy() (bool, time.Time) {
	return t.healthy, t.lastHealthCheck
}

func (t *ListenBrainzTarget) TargetType() TargetType {
	return t.targetType
}

func (t *ListenBrainzTarget) SubmitPlayingTrack(track *common.Track) error {
	req := t.toPlayingNowSubmitListens(track)

	err := t.client.SubmitListens(req)
	if err != nil {
		t.healthy = false
		t.lastHealthCheck = time.Now().UTC()
		return err
	}

	t.healthy = true
	t.lastHealthCheck = time.Now().UTC()
	return nil
}

func (t *ListenBrainzTarget) SubmitPlayedTrack(trackedTrack *common.TrackedTrack) error {
	req := t.toSingleSubmitListens(trackedTrack)

	err := t.client.SubmitListens(req)
	if err != nil {
		t.healthy = false
		t.lastHealthCheck = time.Now().UTC()
		return err
	}

	t.healthy = true
	t.lastHealthCheck = time.Now().UTC()
	return nil
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
					AdditionalInfo: listenbrainz.AdditionalInfo{
						ArtistNames:             track.Artists,
						SubmissionClient:        "Scrobblet",
						SubmissionClientVersion: t.scrobbletVersion,
						TrackNumber:             strconv.Itoa(track.TrackNumber),
					},
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
					AdditionalInfo: listenbrainz.AdditionalInfo{
						ArtistNames:             trackedTrack.Track.Artists,
						SubmissionClient:        "Scrobblet",
						SubmissionClientVersion: t.scrobbletVersion,
						TrackNumber:             strconv.Itoa(trackedTrack.Track.TrackNumber),
					},
				},
			},
		},
	}
}
