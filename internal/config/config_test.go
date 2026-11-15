package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "scrobblet.yaml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	return path
}

func loadConfigFromString(t *testing.T, content string) (*Config, error) {
	t.Helper()

	t.Setenv(envConfig, writeConfigFile(t, content))

	return LoadConfig()
}

const minimalConfig = `
source:
  type: Spotify
  clientId: spotify-id
  clientSecret: spotify-secret
  redirectUrl: http://127.0.0.1:7276/api/spotify/callback
targets:
  - type: CSV
`

func TestLoadConfig_MinimalConfigAppliesDefaults(t *testing.T) {
	cfg, err := loadConfigFromString(t, minimalConfig)
	if err != nil {
		t.Fatalf("LoadConfig() err = %v", err)
	}

	checks := map[string]struct {
		got  any
		want any
	}{
		"Port":                    {cfg.Port, "7276"},
		"DataPath":                {cfg.DataPath, defaultDataPath},
		"LogLevel":                {cfg.LogLevel, defaultLogLevel},
		"RateLimitRate":           {cfg.RateLimitRate, defaultRateLimitRate},
		"RateLimitBurst":          {cfg.RateLimitBurst, defaultRateLimitBurst},
		"Source.Type":             {cfg.Source.Type, sources.SourceSpotify},
		"Source.Spotify.ClientID": {cfg.Source.Spotify.ClientID, "spotify-id"},
		"Source.Spotify.DataPath": {cfg.Source.Spotify.DataPath, defaultDataPath},
	}
	for field, check := range checks {
		if check.got != check.want {
			t.Errorf("%s = %v, want %v", field, check.got, check.want)
		}
	}

	if len(cfg.Targets) != 1 {
		t.Fatalf("len(Targets) = %d, want 1", len(cfg.Targets))
	}

	assertCSVSpec(t, cfg.Targets[0], "csv", filepath.Join(defaultDataPath, "scrobbles.csv"))
}

func TestLoadConfig_GeneralSettingsOverrideDefaults(t *testing.T) {
	cfg, err := loadConfigFromString(t, `
port: 8080
dataPath: /data
logLevel: DEBUG
rateLimit:
  rate: 5
  burst: 50
source:
  type: WiiM
  url: https://192.168.0.100
targets:
  - type: CSV
`)
	if err != nil {
		t.Fatalf("LoadConfig() err = %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want %q", cfg.Port, "8080")
	}
	if cfg.DataPath != "/data" {
		t.Errorf("DataPath = %q, want %q", cfg.DataPath, "/data")
	}
	if cfg.LogLevel != "DEBUG" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "DEBUG")
	}
	if cfg.RateLimitRate != 5 {
		t.Errorf("RateLimitRate = %d, want 5", cfg.RateLimitRate)
	}
	if cfg.RateLimitBurst != 50 {
		t.Errorf("RateLimitBurst = %d, want 50", cfg.RateLimitBurst)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	t.Setenv(envConfig, filepath.Join(t.TempDir(), "does-not-exist.yaml"))

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("LoadConfig() err = nil, want error")
	}
	if !strings.Contains(err.Error(), "does-not-exist.yaml") {
		t.Errorf("error %q does not contain the config file path", err.Error())
	}
}

func TestLoadConfig_EnvVarExpansion(t *testing.T) {
	t.Setenv("SCROBBLET_TEST_TOKEN", "secret-token")

	cfg, err := loadConfigFromString(t, `
source:
  type: WiiM
  url: https://192.168.0.100
targets:
  - type: ListenBrainz
    token: ${SCROBBLET_TEST_TOKEN}
`)
	if err != nil {
		t.Fatalf("LoadConfig() err = %v", err)
	}

	if cfg.Targets[0].ListenBrainz.Token != "secret-token" {
		t.Errorf("Targets[0].ListenBrainz.Token = %q, want %q", cfg.Targets[0].ListenBrainz.Token, "secret-token")
	}
}

func TestLoadConfig_UnsetEnvVarFailsRequiredFieldValidation(t *testing.T) {
	_, err := loadConfigFromString(t, `
source:
  type: WiiM
  url: https://192.168.0.100
targets:
  - type: ListenBrainz
    token: ${SCROBBLET_UNSET_TOKEN}
`)
	if err == nil {
		t.Fatal("LoadConfig() err = nil, want error")
	}
	if !strings.Contains(err.Error(), "token is required") {
		t.Errorf("error %q does not name the required field", err.Error())
	}
}

func TestLoadConfig_UnknownKeysRejected(t *testing.T) {
	tests := []struct {
		name   string
		config string
	}{
		{
			name: "unknown root key",
			config: `
ports: 7276
source:
  type: WiiM
  url: https://192.168.0.100
targets:
  - type: CSV
`,
		},
		{
			name: "unknown source key",
			config: `
source:
  type: WiiM
  url: https://192.168.0.100
  token: abc
targets:
  - type: CSV
`,
		},
		{
			name: "unknown target key",
			config: `
source:
  type: WiiM
  url: https://192.168.0.100
targets:
  - type: ListenBrainz
    tokens: abc
`,
		},
		{
			name: "key of another target type",
			config: `
source:
  type: WiiM
  url: https://192.168.0.100
targets:
  - type: CSV
    apiKey: abc
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadConfigFromString(t, tt.config)
			if err == nil {
				t.Fatal("LoadConfig() err = nil, want error")
			}
		})
	}
}

func TestLoadConfig_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		wantErrText string
	}{
		{
			name: "missing source",
			config: `
targets:
  - type: CSV
`,
			wantErrText: "source is required",
		},
		{
			name: "missing targets",
			config: `
source:
  type: WiiM
  url: https://192.168.0.100
`,
			wantErrText: "at least one target is required",
		},
		{
			name: "unknown source type",
			config: `
source:
  type: AppleMusic
targets:
  - type: CSV
`,
			wantErrText: "invalid source type",
		},
		{
			name: "unknown target type",
			config: `
source:
  type: WiiM
  url: https://192.168.0.100
targets:
  - type: Plex
`,
			wantErrText: "targets[0]: invalid target type",
		},
		{
			name: "duplicate type and name",
			config: `
source:
  type: WiiM
  url: https://192.168.0.100
targets:
  - type: ListenBrainz
    token: a
  - type: ListenBrainz
    token: b
`,
			wantErrText: "duplicate target",
		},
		{
			name: "invalid target name",
			config: `
source:
  type: WiiM
  url: https://192.168.0.100
targets:
  - type: ListenBrainz
    name: Work
    token: a
`,
			wantErrText: "invalid name",
		},
		{
			name: "missing spotify clientId",
			config: `
source:
  type: Spotify
  clientSecret: secret
  redirectUrl: http://127.0.0.1:7276/api/spotify/callback
targets:
  - type: CSV
`,
			wantErrText: "source (Spotify): clientId is required",
		},
		{
			name: "invalid spotify redirectUrl path",
			config: `
source:
  type: Spotify
  clientId: id
  clientSecret: secret
  redirectUrl: http://127.0.0.1:7276/api/wrong/callback
targets:
  - type: CSV
`,
			wantErrText: "redirectUrl",
		},
		{
			name: "wiim url must be a private IP",
			config: `
source:
  type: WiiM
  url: https://8.8.8.8
targets:
  - type: CSV
`,
			wantErrText: "local network IP address",
		},
		{
			name: "koito requires url",
			config: `
source:
  type: WiiM
  url: https://192.168.0.100
targets:
  - type: Koito
    token: a
`,
			wantErrText: "url is required",
		},
		{
			name: "lastfm default redirectUrl path",
			config: `
source:
  type: WiiM
  url: https://192.168.0.100
targets:
  - type: LastFm
    apiKey: key
    sharedSecret: secret
    redirectUrl: http://127.0.0.1:7276/api/lastfm/work/callback
`,
			wantErrText: "Path must be \"/api/lastfm/callback\"",
		},
		{
			name: "named lastfm redirectUrl path",
			config: `
source:
  type: WiiM
  url: https://192.168.0.100
targets:
  - type: LastFm
    name: work
    apiKey: key
    sharedSecret: secret
    redirectUrl: http://127.0.0.1:7276/api/lastfm/callback
`,
			wantErrText: "Path must be \"/api/lastfm/work/callback\"",
		},
		{
			name: "csv filePath conflict",
			config: `
source:
  type: WiiM
  url: https://192.168.0.100
targets:
  - type: CSV
    filePath: /data/scrobbles.csv
  - type: CSV
    name: archive
    filePath: /data/scrobbles.csv
`,
			wantErrText: "CSV targets must write to different files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadConfigFromString(t, tt.config)
			if err == nil {
				t.Fatal("LoadConfig() err = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantErrText) {
				t.Errorf("error %q does not contain %q", err.Error(), tt.wantErrText)
			}
		})
	}
}

// Golden test: the full example from the configuration docs loads into the
// expected construction specs.
func TestLoadConfig_FullExample(t *testing.T) {
	t.Setenv("SPOTIFY_CLIENT_ID", "spotify-id")
	t.Setenv("SPOTIFY_CLIENT_SECRET", "spotify-secret")
	t.Setenv("LISTENBRAINZ_TOKEN", "lb-token")
	t.Setenv("LISTENBRAINZ_WORK_TOKEN", "lb-work-token")
	t.Setenv("MALOJA_TOKEN", "maloja-token")
	t.Setenv("KOITO_TOKEN", "koito-token")
	t.Setenv("LASTFM_API_KEY", "lastfm-key")
	t.Setenv("LASTFM_SHARED_SECRET", "lastfm-secret")
	t.Setenv("LASTFM_2_API_KEY", "lastfm-2-key")
	t.Setenv("LASTFM_2_SHARED_SECRET", "lastfm-2-secret")

	cfg, err := loadConfigFromString(t, `
port: 7276
dataPath: /data
logLevel: INFO
rateLimit:
  rate: 10
  burst: 100

source:
  type: Spotify
  clientId: ${SPOTIFY_CLIENT_ID}
  clientSecret: ${SPOTIFY_CLIENT_SECRET}
  redirectUrl: http://127.0.0.1:7276/api/spotify/callback

targets:
  - type: ListenBrainz
    token: ${LISTENBRAINZ_TOKEN}

  - type: ListenBrainz
    name: work
    url: https://listenbrainz.example.com
    token: ${LISTENBRAINZ_WORK_TOKEN}

  - type: Maloja
    url: http://maloja:42010
    token: ${MALOJA_TOKEN}

  - type: Koito
    url: http://koito:4110
    token: ${KOITO_TOKEN}

  - type: LastFm
    apiKey: ${LASTFM_API_KEY}
    sharedSecret: ${LASTFM_SHARED_SECRET}
    redirectUrl: http://127.0.0.1:7276/api/lastfm/callback

  - type: LastFm
    name: second-account
    apiKey: ${LASTFM_2_API_KEY}
    sharedSecret: ${LASTFM_2_SHARED_SECRET}
    redirectUrl: http://127.0.0.1:7276/api/lastfm/second-account/callback

  - type: CSV

  - type: CSV
    name: archive
`)
	if err != nil {
		t.Fatalf("LoadConfig() err = %v", err)
	}

	if cfg.Source.Type != sources.SourceSpotify {
		t.Errorf("Source.Type = %q, want %q", cfg.Source.Type, sources.SourceSpotify)
	}
	if cfg.Source.Spotify.ClientID != "spotify-id" {
		t.Errorf("Source.Spotify.ClientID = %q, want %q", cfg.Source.Spotify.ClientID, "spotify-id")
	}

	if len(cfg.Targets) != 8 {
		t.Fatalf("len(Targets) = %d, want 8", len(cfg.Targets))
	}

	assertListenBrainzSpec(t, cfg.Targets[0], targets.TargetListenBrainz, "listenbrainz", "https://api.listenbrainz.org", "lb-token")
	assertListenBrainzSpec(t, cfg.Targets[1], targets.TargetListenBrainz, "work", "https://listenbrainz.example.com", "lb-work-token")
	assertListenBrainzSpec(t, cfg.Targets[2], targets.TargetMaloja, "maloja", "http://maloja:42010/apis/listenbrainz", "maloja-token")
	assertListenBrainzSpec(t, cfg.Targets[3], targets.TargetKoito, "koito", "http://koito:4110/apis/listenbrainz", "koito-token")
	assertLastFmSpec(t, cfg.Targets[4], "lastfm", "lastfm-key", "http://127.0.0.1:7276/api/lastfm/callback")
	assertLastFmSpec(t, cfg.Targets[5], "second-account", "lastfm-2-key", "http://127.0.0.1:7276/api/lastfm/second-account/callback")
	assertCSVSpec(t, cfg.Targets[6], "csv", filepath.Join("/data", "scrobbles.csv"))
	assertCSVSpec(t, cfg.Targets[7], "archive", filepath.Join("/data", "scrobbles-archive.csv"))
}

func assertListenBrainzSpec(t *testing.T, spec targets.Spec, targetType targets.TargetType, name, url, token string) {
	t.Helper()

	if spec.Type != targetType {
		t.Errorf("spec.Type = %q, want %q", spec.Type, targetType)
	}
	if spec.Name != name {
		t.Errorf("spec.Name = %q, want %q", spec.Name, name)
	}
	if spec.ListenBrainz == nil {
		t.Fatalf("spec.ListenBrainz is nil for target %q", name)
	}
	if spec.ListenBrainz.URL != url {
		t.Errorf("spec.ListenBrainz.URL = %q, want %q", spec.ListenBrainz.URL, url)
	}
	if spec.ListenBrainz.Token != token {
		t.Errorf("spec.ListenBrainz.Token = %q, want %q", spec.ListenBrainz.Token, token)
	}
}

func assertLastFmSpec(t *testing.T, spec targets.Spec, name, apiKey, redirectURL string) {
	t.Helper()

	if spec.Type != targets.TargetLastFm {
		t.Errorf("spec.Type = %q, want %q", spec.Type, targets.TargetLastFm)
	}
	if spec.Name != name {
		t.Errorf("spec.Name = %q, want %q", spec.Name, name)
	}
	if spec.LastFm == nil {
		t.Fatalf("spec.LastFm is nil for target %q", name)
	}
	if spec.LastFm.APIKey != apiKey {
		t.Errorf("spec.LastFm.APIKey = %q, want %q", spec.LastFm.APIKey, apiKey)
	}
	if spec.LastFm.RedirectURL != redirectURL {
		t.Errorf("spec.LastFm.RedirectURL = %q, want %q", spec.LastFm.RedirectURL, redirectURL)
	}
	if spec.LastFm.DataPath != "/data" {
		t.Errorf("spec.LastFm.DataPath = %q, want %q", spec.LastFm.DataPath, "/data")
	}
}

func assertCSVSpec(t *testing.T, spec targets.Spec, name, filePath string) {
	t.Helper()

	if spec.Type != targets.TargetCSV {
		t.Errorf("spec.Type = %q, want %q", spec.Type, targets.TargetCSV)
	}
	if spec.Name != name {
		t.Errorf("spec.Name = %q, want %q", spec.Name, name)
	}
	if spec.CSV == nil {
		t.Fatalf("spec.CSV is nil for target %q", name)
	}
	if spec.CSV.FilePath != filePath {
		t.Errorf("spec.CSV.FilePath = %q, want %q", spec.CSV.FilePath, filePath)
	}
}
