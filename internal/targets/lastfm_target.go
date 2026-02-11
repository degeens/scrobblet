package targets

import (
	"time"

	"github.com/degeens/scrobblet/internal/clients/lastfm"
	"github.com/degeens/scrobblet/internal/common"
)

type LastFmTarget struct {
	healthy         bool
	lastHealthCheck time.Time
	client          *lastfm.Client
}

func NewLastFmTarget(client *lastfm.Client) *LastFmTarget {
	return &LastFmTarget{
		healthy:         true,
		lastHealthCheck: time.Now(),
		client:          client,
	}
}

func (t *LastFmTarget) Healthy() (bool, time.Time) {
	return t.healthy, t.lastHealthCheck
}

func (t *LastFmTarget) TargetType() TargetType {
	return TargetLastFm
}

func (t *LastFmTarget) SubmitPlayingTrack(track *common.Track) error {
	req := t.toUpdateNowPlaying(track)

	err := t.client.UpdateNowPlaying(req)
	if err != nil {
		t.healthy = false
		t.lastHealthCheck = time.Now()
		return err
	}

	t.healthy = true
	t.lastHealthCheck = time.Now()
	return nil
}

func (t *LastFmTarget) SubmitPlayedTrack(trackedTrack *common.TrackedTrack) error {
	req := t.toScrobble(trackedTrack)

	err := t.client.Scrobble([]lastfm.ScrobbleRequest{req})
	if err != nil {
		t.healthy = false
		t.lastHealthCheck = time.Now()
		return err
	}

	t.healthy = true
	t.lastHealthCheck = time.Now()
	return nil
}

func (t *LastFmTarget) toUpdateNowPlaying(track *common.Track) *lastfm.UpdateNowPlayingRequest {
	// Last.fm does only support one artist
	artistName := track.Artists[0]

	return &lastfm.UpdateNowPlayingRequest{
		Artist:      artistName,
		Track:       track.Title,
		Album:       track.Album,
		Duration:    int(track.Duration.Seconds()),
		TrackNumber: track.TrackNumber,
	}
}

func (t *LastFmTarget) toScrobble(trackedTrack *common.TrackedTrack) lastfm.ScrobbleRequest {
	// Last.fm does only support one artist
	artistName := trackedTrack.Track.Artists[0]

	return lastfm.ScrobbleRequest{
		Artist:      artistName,
		Track:       trackedTrack.Track.Title,
		Timestamp:   trackedTrack.StartedAt.Unix(),
		Album:       trackedTrack.Track.Album,
		Duration:    int(trackedTrack.Track.Duration.Seconds()),
		TrackNumber: trackedTrack.Track.TrackNumber,
	}
}
