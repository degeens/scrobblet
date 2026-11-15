package targets

import (
	"testing"

	"github.com/degeens/scrobblet/internal/clients/csv"
	"github.com/degeens/scrobblet/internal/clients/lastfm"
	"github.com/degeens/scrobblet/internal/clients/listenbrainz"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		spec       Spec
		wantType   TargetType
		wantName   string
		wantClient any
		wantErr    bool
	}{
		{
			name: "listenbrainz",
			spec: Spec{
				Type:         TargetListenBrainz,
				Name:         "listenbrainz",
				ListenBrainz: &listenbrainz.Config{URL: "https://api.listenbrainz.org", Token: "token"},
			},
			wantType:   TargetListenBrainz,
			wantName:   "listenbrainz",
			wantClient: &listenbrainz.Client{},
		},
		{
			name: "koito",
			spec: Spec{
				Type:         TargetKoito,
				Name:         "koito",
				ListenBrainz: &listenbrainz.Config{URL: "http://koito:4110/apis/listenbrainz", Token: "token"},
			},
			wantType:   TargetKoito,
			wantName:   "koito",
			wantClient: &listenbrainz.Client{},
		},
		{
			name: "maloja",
			spec: Spec{
				Type:         TargetMaloja,
				Name:         "maloja",
				ListenBrainz: &listenbrainz.Config{URL: "http://maloja:42010/apis/listenbrainz", Token: "token"},
			},
			wantType:   TargetMaloja,
			wantName:   "maloja",
			wantClient: &listenbrainz.Client{},
		},
		{
			name: "lastfm",
			spec: Spec{
				Type: TargetLastFm,
				Name: "lastfm",
				LastFm: &lastfm.Config{
					APIKey:       "key",
					SharedSecret: "secret",
					RedirectURL:  "http://127.0.0.1:7276/api/lastfm/callback",
					DataPath:     t.TempDir(),
				},
			},
			wantType:   TargetLastFm,
			wantName:   "lastfm",
			wantClient: &lastfm.Client{},
		},
		{
			name: "csv",
			spec: Spec{
				Type: TargetCSV,
				Name: "csv",
				CSV:  &csv.Config{FilePath: "scrobbles.csv"},
			},
			wantType:   TargetCSV,
			wantName:   "csv",
			wantClient: &csv.Client{},
		},
		{
			name: "unknown type",
			spec: Spec{
				Type: "Unknown",
				Name: "unknown",
			},
			wantErr: true,
		},
		{
			name: "missing type config",
			spec: Spec{
				Type: TargetCSV,
				Name: "csv",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance, err := New(tt.spec, "test")
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if instance.Type != tt.wantType {
				t.Errorf("instance.Type = %q, want %q", instance.Type, tt.wantType)
			}
			if instance.Name != tt.wantName {
				t.Errorf("instance.Name = %q, want %q", instance.Name, tt.wantName)
			}

			switch tt.wantClient.(type) {
			case *listenbrainz.Client:
				if _, ok := instance.Client.(*listenbrainz.Client); !ok {
					t.Errorf("instance.Client = %T, want *listenbrainz.Client", instance.Client)
				}
			case *lastfm.Client:
				if _, ok := instance.Client.(*lastfm.Client); !ok {
					t.Errorf("instance.Client = %T, want *lastfm.Client", instance.Client)
				}
			case *csv.Client:
				if _, ok := instance.Client.(*csv.Client); !ok {
					t.Errorf("instance.Client = %T, want *csv.Client", instance.Client)
				}
			}

			if instance.Target.Name() != tt.wantName {
				t.Errorf("instance.Target.Name() = %q, want %q", instance.Target.Name(), tt.wantName)
			}
			if instance.Target.TargetType() != tt.wantType {
				t.Errorf("instance.Target.TargetType() = %q, want %q", instance.Target.TargetType(), tt.wantType)
			}
		})
	}
}

func TestNewMultiple_SameTypeInstances(t *testing.T) {
	specs := []Spec{
		{
			Type:         TargetListenBrainz,
			Name:         "listenbrainz",
			ListenBrainz: &listenbrainz.Config{Token: "token-a"},
		},
		{
			Type:         TargetListenBrainz,
			Name:         "work",
			ListenBrainz: &listenbrainz.Config{URL: "https://listenbrainz.example.com", Token: "token-b"},
		},
	}

	instances, err := NewMultiple(specs, "test")
	if err != nil {
		t.Fatalf("NewMultiple() err = %v", err)
	}

	if len(instances) != 2 {
		t.Fatalf("len(instances) = %d, want 2", len(instances))
	}

	if instances[0].Name != "listenbrainz" || instances[0].Type != TargetListenBrainz {
		t.Errorf("instances[0] = (%q, %q), want (%q, %q)", instances[0].Name, instances[0].Type, "listenbrainz", TargetListenBrainz)
	}
	if instances[1].Name != "work" || instances[1].Type != TargetListenBrainz {
		t.Errorf("instances[1] = (%q, %q), want (%q, %q)", instances[1].Name, instances[1].Type, "work", TargetListenBrainz)
	}
}

func TestNewMultiple_ErrorAborts(t *testing.T) {
	specs := []Spec{
		{
			Type: TargetCSV,
			Name: "csv",
			CSV:  &csv.Config{FilePath: "scrobbles.csv"},
		},
		{
			Type: "Unknown",
			Name: "unknown",
		},
	}

	_, err := NewMultiple(specs, "test")
	if err == nil {
		t.Fatal("NewMultiple() err = nil, want error")
	}
}
