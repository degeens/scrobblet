package targets

import (
	"fmt"

	"github.com/degeens/scrobblet/internal/clients"
	"github.com/degeens/scrobblet/internal/clients/koito"
	"github.com/degeens/scrobblet/internal/clients/listenbrainz"
	"github.com/degeens/scrobblet/internal/common"
)

type TargetType string

const (
	TargetKoito        TargetType = "Koito"
	TargetListenBrainz TargetType = "ListenBrainz"
)

type Target interface {
	SubmitPlayingTrack(track *common.Track) error
	SubmitPlayedTrack(trackedTrack *common.TrackedTrack) error
}

func New(targetType TargetType, clientsConfig clients.Config) (any, Target, error) {
	switch targetType {
	case TargetKoito:
		client := koito.NewClient(clientsConfig.Koito.URL, clientsConfig.Koito.Token)
		return client, NewKoitoTarget(client), nil
	case TargetListenBrainz:
		client := listenbrainz.NewClient(clientsConfig.ListenBrainz.Token)
		return client, NewListenBrainzTarget(client), nil
	default:
		return nil, nil, fmt.Errorf("Unknown target type: %s", targetType)
	}
}
