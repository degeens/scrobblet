package targets

import (
	"fmt"
	"time"

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
	Healthy() (bool, time.Time)
	Name() string
	TargetType() TargetType
	SubmitPlayingTrack(track *common.Track) error
	SubmitPlayedTrack(trackedTrack *common.TrackedTrack) error
}

// Spec describes how to construct a target instance. The config field
// matching Type must be set.
type Spec struct {
	Type         TargetType
	Name         string
	ListenBrainz *listenbrainz.Config
	LastFm       *lastfm.Config
	CSV          *csv.Config
}

// Instance is a constructed target together with its identity and client.
type Instance struct {
	Name   string
	Type   TargetType
	Client any
	Target Target
}

func New(spec Spec, scrobbletVersion string) (Instance, error) {
	switch spec.Type {
	case TargetKoito, TargetMaloja, TargetListenBrainz:
		// Koito and Maloja use the ListenBrainz-compatible API with a custom base URL
		if spec.ListenBrainz == nil {
			return Instance{}, fmt.Errorf("missing ListenBrainz configuration for target type: %s", spec.Type)
		}
		client := listenbrainz.NewClient(spec.ListenBrainz.Token, spec.ListenBrainz.URL)
		return Instance{
			Name:   spec.Name,
			Type:   spec.Type,
			Client: client,
			Target: NewListenBrainzTarget(spec.Type, spec.Name, client, scrobbletVersion),
		}, nil
	case TargetLastFm:
		if spec.LastFm == nil {
			return Instance{}, fmt.Errorf("missing Last.fm configuration for target type: %s", spec.Type)
		}
		client, err := lastfm.NewClient(spec.Name, spec.LastFm.APIKey, spec.LastFm.SharedSecret, spec.LastFm.RedirectURL, spec.LastFm.DataPath, scrobbletVersion)
		if err != nil {
			return Instance{}, err
		}
		return Instance{
			Name:   spec.Name,
			Type:   spec.Type,
			Client: client,
			Target: NewLastFmTarget(spec.Name, client),
		}, nil
	case TargetCSV:
		if spec.CSV == nil {
			return Instance{}, fmt.Errorf("missing CSV configuration for target type: %s", spec.Type)
		}
		client := csv.NewClient(spec.CSV.FilePath)
		return Instance{
			Name:   spec.Name,
			Type:   spec.Type,
			Client: client,
			Target: NewCSVTarget(spec.Name, client),
		}, nil
	default:
		return Instance{}, fmt.Errorf("unknown target type: %s", spec.Type)
	}
}

func NewMultiple(specs []Spec, scrobbletVersion string) ([]Instance, error) {
	instances := make([]Instance, 0, len(specs))

	for _, spec := range specs {
		instance, err := New(spec, scrobbletVersion)
		if err != nil {
			return nil, err
		}

		instances = append(instances, instance)
	}

	return instances, nil
}
