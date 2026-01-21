package scrobbler

import (
	"testing"
	"time"
)

func TestShouldScrobble(t *testing.T) {
	tests := []struct {
		name            string
		trackedDuration time.Duration
		trackDuration   time.Duration
		want            bool
	}{
		{
			name:            "track <= 30s",
			trackedDuration: 30 * time.Second,
			trackDuration:   30 * time.Second,
			want:            false,
		},
		{
			name:            "track > 30s, played < 50% and < 4min",
			trackedDuration: 2 * time.Minute,
			trackDuration:   6 * time.Minute,
			want:            false,
		},
		{
			name:            "track > 30s, played >= 50% and < 4min",
			trackedDuration: 3 * time.Minute,
			trackDuration:   6 * time.Minute,
			want:            true,
		},
		{
			name:            "track > 30s, played < 50% and >= 4min",
			trackedDuration: 4 * time.Minute,
			trackDuration:   9 * time.Minute,
			want:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldScrobble(tt.trackedDuration, tt.trackDuration)
			if got != tt.want {
				t.Errorf("ShouldScrobble(%v, %v) = %v, want %v",
					tt.trackedDuration, tt.trackDuration, got, tt.want)
			}
		})
	}
}
