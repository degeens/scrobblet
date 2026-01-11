package scrobbler

import "time"

func ShouldScrobble(trackedDuration, trackDuration time.Duration) bool {
	// The track must be longer than 30 seconds
	thirtySeconds := time.Duration(30) * time.Second
	if trackDuration <= thirtySeconds {
		return false
	}

	// And the track has been played for at least half its duration, or for 4 minutes
	fourMinutes := time.Duration(4) * time.Minute
	halfTrack := trackDuration / 2

	return trackedDuration >= halfTrack || trackedDuration >= fourMinutes
}
