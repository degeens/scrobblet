package targets

import (
	"fmt"

	"github.com/degeens/scrobblet/internal/clients"
	"github.com/degeens/scrobblet/internal/clients/csv"
	"github.com/degeens/scrobblet/internal/clients/lastfm"
	"github.com/degeens/scrobblet/internal/clients/listenbrainz"
	"github.com/degeens/scrobblet/internal/common"
)

type TargetType string

const (
	TargetKoito        TargetType = "Koito"
	TargetMaloja       TargetType = "Maloja"
	TargetListenBrainz TargetType = "ListenBrainz"
	TargetLastFm       TargetType = "LastFm"
	TargetCSV          TargetType = "CSV"
)

type Target interface {
	SubmitPlayingTrack(track *common.Track) error
	SubmitPlayedTrack(trackedTrack *common.TrackedTrack) error
}

func New(targetType TargetType, clientsConfig clients.Config, scrobbletVersion string) (any, Target, error) {
	switch targetType {
	case TargetKoito:
		// Koito uses ListenBrainz-compatible API with custom base URL
		client := listenbrainz.NewClient(clientsConfig.ListenBrainz.Token, clientsConfig.ListenBrainz.URL)
		return client, NewListenBrainzTarget(client, scrobbletVersion), nil
	case TargetMaloja:
		// Maloja uses ListenBrainz-compatible API with custom base URL
		client := listenbrainz.NewClient(clientsConfig.ListenBrainz.Token, clientsConfig.ListenBrainz.URL)
		return client, NewListenBrainzTarget(client, scrobbletVersion), nil
	case TargetListenBrainz:
		client := listenbrainz.NewClient(clientsConfig.ListenBrainz.Token, clientsConfig.ListenBrainz.URL)
		return client, NewListenBrainzTarget(client, scrobbletVersion), nil
	case TargetLastFm:
		client, err := lastfm.NewClient(clientsConfig.LastFm.APIKey, clientsConfig.LastFm.SharedSecret, clientsConfig.LastFm.RedirectURL, clientsConfig.LastFm.DataPath, scrobbletVersion)
		if err != nil {
			return nil, nil, err
		}
		return client, NewLastFmTarget(client), nil
	case TargetCSV:
		client := csv.NewClient(clientsConfig.CSV.FilePath)
		return client, NewCSVTarget(client), nil
	default:
		return nil, nil, fmt.Errorf("unknown target type: %s", targetType)
	}
}
