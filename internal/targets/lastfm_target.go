package targets

import (
	"github.com/degeens/scrobblet/internal/clients/lastfm"
	"github.com/degeens/scrobblet/internal/common"
)

type LastFmTarget struct {
	client *lastfm.Client
}

func NewLastFmTarget(client *lastfm.Client) *LastFmTarget {
	return &LastFmTarget{
		client: client,
	}
}

func (t *LastFmTarget) SubmitPlayingTrack(track *common.Track) error {
	req := t.toUpdateNowPlaying(track)

	return t.client.UpdateNowPlaying(req)
}

func (t *LastFmTarget) SubmitPlayedTrack(trackedTrack *common.TrackedTrack) error {
	req := t.toScrobble(trackedTrack)

	return t.client.Scrobble([]lastfm.ScrobbleRequest{req})
}

func (t *LastFmTarget) toUpdateNowPlaying(track *common.Track) *lastfm.UpdateNowPlayingRequest {
	// Last.fm does only support one artist
	artistName := track.Artists[0]

	return &lastfm.UpdateNowPlayingRequest{
		Artist:   artistName,
		Track:    track.Title,
		Album:    track.Album,
		Duration: int(track.Duration.Seconds()),
	}
}

func (t *LastFmTarget) toScrobble(trackedTrack *common.TrackedTrack) lastfm.ScrobbleRequest {
	// Last.fm does only support one artist
	artistName := trackedTrack.Track.Artists[0]

	return lastfm.ScrobbleRequest{
		Artist:    artistName,
		Track:     trackedTrack.Track.Title,
		Timestamp: trackedTrack.StartedAt.Unix(),
		Album:     trackedTrack.Track.Album,
		Duration:  int(trackedTrack.Track.Duration.Seconds()),
	}
}
