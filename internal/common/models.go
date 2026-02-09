package common

import (
	"strings"
	"time"
)

type Track struct {
	Artists     []string
	Title       string
	Album       string
	TrackNumber int
	Duration    time.Duration
}

func (t *Track) Equals(other *Track) bool {
	if t.Title != other.Title || t.Album != other.Album || t.TrackNumber != other.TrackNumber || t.Duration != other.Duration {
		return false
	}

	if len(t.Artists) != len(other.Artists) {
		return false
	}

	for i, artist := range t.Artists {
		if artist != other.Artists[i] {
			return false
		}
	}

	return true
}

func (t *Track) SlogArgs() []any {
	return []any{"Artist", strings.Join(t.Artists, ", "), "Title", t.Title, "Album", t.Album}
}

type TrackedTrack struct {
	Track         *Track
	LastPosition  time.Duration
	Duration      time.Duration
	StartedAt     time.Time
	LastUpdatedAt time.Time
}

func (t *TrackedTrack) SlogArgs() []any {
	return append(t.Track.SlogArgs(), "TrackDuration", t.Track.Duration/time.Second, "TrackedDuration", t.Duration/time.Second)
}
