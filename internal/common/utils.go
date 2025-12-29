package common

import "time"

func HasReachedScrobbleThreshold(trackedDuration, trackDuration time.Duration) bool {
	return hasReachedAbsoluteScrobbleThreshold(trackedDuration) || hasReachedRelativeScrobbleThreshold(trackedDuration, trackDuration)
}

func hasReachedAbsoluteScrobbleThreshold(trackedDuration time.Duration) bool {
	threeMinutes := time.Duration(3) * time.Minute

	return trackedDuration >= threeMinutes
}

func hasReachedRelativeScrobbleThreshold(trackedDuration, trackDuration time.Duration) bool {
	halfTrack := trackDuration / 2

	return trackedDuration >= halfTrack
}
