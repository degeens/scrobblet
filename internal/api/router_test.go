package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/degeens/scrobblet/internal/clients/lastfm"
	"github.com/degeens/scrobblet/internal/clients/listenbrainz"
	"github.com/degeens/scrobblet/internal/clients/wiim"
	"github.com/degeens/scrobblet/internal/config"
	"github.com/degeens/scrobblet/internal/metrics"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

func newTestHandler(t *testing.T, targetSpecs []targets.Spec) http.Handler {
	t.Helper()

	sourceClient, source, err := sources.New(sources.Spec{
		Type: sources.SourceWiiM,
		WiiM: wiimConfig(),
	})
	if err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	instances, err := targets.NewMultiple(targetSpecs, "test")
	if err != nil {
		t.Fatalf("failed to create targets: %v", err)
	}

	cfg := &config.Config{RateLimitRate: 100, RateLimitBurst: 100}

	return NewHandler(source, instances, sourceClient, cfg, metrics.New())
}

func wiimConfig() *wiim.Config {
	return &wiim.Config{BaseURL: "https://192.168.0.100"}
}

func lastFmSpec(t *testing.T, name, redirectURL string) targets.Spec {
	t.Helper()

	return targets.Spec{
		Type: targets.TargetLastFm,
		Name: name,
		LastFm: &lastfm.Config{
			APIKey:       "key",
			SharedSecret: "secret",
			RedirectURL:  redirectURL,
			DataPath:     t.TempDir(),
		},
	}
}

func requestStatus(t *testing.T, handler http.Handler, path string) int {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, path, nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	return recorder.Code
}

func TestRouter_LastFmRoutes(t *testing.T) {
	handler := newTestHandler(t, []targets.Spec{
		lastFmSpec(t, "lastfm", "http://127.0.0.1:7276/api/lastfm/callback"),
		lastFmSpec(t, "work", "http://127.0.0.1:7276/api/lastfm/work/callback"),
	})

	tests := []struct {
		path       string
		wantStatus int
	}{
		{"/api/lastfm/login", http.StatusFound},
		{"/api/lastfm/work/login", http.StatusFound},
		// Callbacks without a token parameter are rejected but routed
		{"/api/lastfm/callback", http.StatusBadRequest},
		{"/api/lastfm/work/callback", http.StatusBadRequest},
		// Routes of unknown instances are not registered
		{"/api/lastfm/unknown/login", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if status := requestStatus(t, handler, tt.path); status != tt.wantStatus {
				t.Errorf("GET %s status = %d, want %d", tt.path, status, tt.wantStatus)
			}
		})
	}
}

func TestRouter_HealthIncludesTargetNames(t *testing.T) {
	handler := newTestHandler(t, []targets.Spec{
		{
			Type:         targets.TargetListenBrainz,
			Name:         "listenbrainz",
			ListenBrainz: &listenbrainz.Config{Token: "token-a"},
		},
		{
			Type:         targets.TargetListenBrainz,
			Name:         "work",
			ListenBrainz: &listenbrainz.Config{URL: "https://listenbrainz.example.com", Token: "token-b"},
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	var response struct {
		Source struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"source"`
		Targets []struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"targets"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode health response: %v", err)
	}

	if len(response.Targets) != 2 {
		t.Fatalf("len(targets) = %d, want 2", len(response.Targets))
	}
	if response.Targets[0].Name != "listenbrainz" {
		t.Errorf("targets[0].name = %q, want %q", response.Targets[0].Name, "listenbrainz")
	}
	if response.Targets[1].Name != "work" {
		t.Errorf("targets[1].name = %q, want %q", response.Targets[1].Name, "work")
	}
	if response.Source.Name != "" {
		t.Errorf("source.name = %q, want empty", response.Source.Name)
	}
	if strings.Contains(recorder.Body.String(), `"source":{"type":"WiiM","name"`) {
		t.Error("source should not have a name field")
	}
}
