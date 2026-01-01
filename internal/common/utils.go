package common

import "time"

func HasReachedScrobbleThreshold(trackedDuration, trackDuration time.Duration) bool {
	return hasReachedAbsoluteScrobbleThreshold(trackedDuration) || hasReachedRelativeScrobbleThreshold(trackedDuration, trackDuration)
}

func hasReachedAbsoluteScrobbleThreshold(trackedDuration time.Duration) bool {
	fourMinutes := time.Duration(4) * time.Minute

	return trackedDuration >= fourMinutes
}

func hasReachedRelativeScrobbleThreshold(trackedDuration, trackDuration time.Duration) bool {
	halfTrack := trackDuration / 2

	return trackedDuration >= halfTrack
}
