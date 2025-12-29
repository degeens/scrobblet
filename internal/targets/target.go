package targets

import "github.com/degeens/scrobblet/internal/common"

type Target interface {
	SubmitTrack(track *common.TrackedTrack) error
}
