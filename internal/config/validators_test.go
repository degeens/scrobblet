package config

import (
	"testing"

	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

func TestValidateSourceType(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    sources.SourceType
		wantErr bool
	}{
		{
			name:    "valid spotify lowercase",
			source:  "spotify",
			want:    sources.SourceSpotify,
			wantErr: false,
		},
		{
			name:    "valid spotify titlecase",
			source:  "Spotify",
			want:    sources.SourceSpotify,
			wantErr: false,
		},
		{
			name:    "valid with surrounding spaces",
			source:  " Spotify ",
			want:    sources.SourceSpotify,
			wantErr: false,
		},
		{
			name:    "valid wiim lowercase",
			source:  "wiim",
			want:    sources.SourceWiiM,
			wantErr: false,
		},
		{
			name:    "valid wiim mixed case",
			source:  "WiiM",
			want:    sources.SourceWiiM,
			wantErr: false,
		},
		{
			name:    "unknown returns error",
			source:  "unknown",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty string returns error",
			source:  "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateSourceType(tt.source)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateSourceType(%q) err = %v, wantErr %v", tt.source, err != nil, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("validateSourceType(%q) = %q, want %q", tt.source, got, tt.want)
			}
		})
	}
}

func TestValidateTargetType(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		want    targets.TargetType
		wantErr bool
	}{
		{
			name:    "valid koito lowercase",
			target:  "koito",
			want:    targets.TargetKoito,
			wantErr: false,
		},
		{
			name:    "valid csv titlecase",
			target:  "CSV",
			want:    targets.TargetCSV,
			wantErr: false,
		},
		{
			name:    "valid with surrounding spaces",
			target:  " ListenBrainz ",
			want:    targets.TargetListenBrainz,
			wantErr: false,
		},
		{
			name:    "valid maloja",
			target:  "Maloja",
			want:    targets.TargetMaloja,
			wantErr: false,
		},
		{
			name:    "valid lastfm",
			target:  "LastFm",
			want:    targets.TargetLastFm,
			wantErr: false,
		},
		{
			name:    "unknown returns error",
			target:  "unknown",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty string returns error",
			target:  "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateTargetType(tt.target)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateTargetType(%q) err = %v, wantErr %v", tt.target, err != nil, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("validateTargetType(%q) = %q, want %q", tt.target, got, tt.want)
			}
		})
	}
}

func TestValidateTargetName(t *testing.T) {
	tests := []struct {
		name       string
		targetName string
		wantErr    bool
	}{
		{name: "single letter", targetName: "a", wantErr: false},
		{name: "single digit", targetName: "1", wantErr: false},
		{name: "lowercase with digits", targetName: "work2", wantErr: false},
		{name: "with dashes", targetName: "second-account", wantErr: false},
		{name: "uppercase rejected", targetName: "Work", wantErr: true},
		{name: "leading dash rejected", targetName: "-work", wantErr: true},
		{name: "underscore rejected", targetName: "second_account", wantErr: true},
		{name: "space rejected", targetName: "second account", wantErr: true},
		{name: "empty rejected", targetName: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTargetName(tt.targetName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTargetName(%q) err = %v, wantErr %v", tt.targetName, err != nil, tt.wantErr)
			}
		})
	}
}

func TestValidateRedirectURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		validPath   string
		wantErr     bool
		wantErrText string
	}{
		{
			name:        "valid http scheme, host, and path",
			url:         "http://example.com/spotify/callback",
			validPath:   "/spotify/callback",
			wantErr:     false,
			wantErrText: "",
		},
		{
			name:        "valid https scheme, host, and path",
			url:         "https://example.com/lastfm/callback",
			validPath:   "/lastfm/callback",
			wantErr:     false,
			wantErrText: "",
		},
		{
			name:        "invalid scheme",
			url:         "ftp://example.com/spotify/callback",
			validPath:   "/spotify/callback",
			wantErr:     true,
			wantErrText: "invalid URL scheme: \"ftp\". Scheme must be http or https",
		},
		{
			name:        "empty host",
			url:         "http:///spotify/callback",
			validPath:   "/spotify/callback",
			wantErr:     true,
			wantErrText: "invalid URL: host must not be empty",
		},
		{
			name:        "invalid path",
			url:         "http://example.com/other",
			validPath:   "/spotify/callback",
			wantErr:     true,
			wantErrText: "invalid URL path: \"/other\". Path must be \"/spotify/callback\"",
		},
		{
			name:        "empty url",
			url:         "",
			validPath:   "/spotify/callback",
			wantErr:     true,
			wantErrText: "invalid URL scheme: \"\". Scheme must be http or https",
		},
		{
			name:        "missing scheme",
			url:         "example.com/spotify/callback",
			validPath:   "/spotify/callback",
			wantErr:     true,
			wantErrText: "invalid URL scheme: \"\". Scheme must be http or https",
		},
		{
			name:        "missing path",
			url:         "http://example.com",
			validPath:   "/spotify/callback",
			wantErr:     true,
			wantErrText: "invalid URL path: \"\". Path must be \"/spotify/callback\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRedirectURL(tt.url, tt.validPath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateRedirectURL(%q, %q) err = %v, wantErr %v", tt.url, tt.validPath, err != nil, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.wantErrText {
				t.Errorf("validateRedirectURL(%q, %q) errText = %q, wantErrText %q", tt.url, tt.validPath, err.Error(), tt.wantErrText)
			}
		})
	}
}

func TestValidateWiiMURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantErr     bool
		wantErrText string
	}{
		{
			name:        "valid 192.168.x range",
			url:         "https://192.168.0.100",
			wantErr:     false,
			wantErrText: "",
		},
		{
			name:        "valid 192.168.x range with port",
			url:         "https://192.168.0.100:8443",
			wantErr:     false,
			wantErrText: "",
		},
		{
			name:        "valid 10.x range",
			url:         "http://10.0.0.5",
			wantErr:     false,
			wantErrText: "",
		},
		{
			name:        "valid 172.16.x range",
			url:         "http://172.16.5.4",
			wantErr:     false,
			wantErrText: "",
		},
		{
			name:        "public IP",
			url:         "https://8.8.8.8",
			wantErr:     true,
			wantErrText: "invalid URL host: \"8.8.8.8\". Host must be a local network IP address",
		},
		{
			name:        "hostname",
			url:         "https://wiim.local",
			wantErr:     true,
			wantErrText: "invalid URL host: \"wiim.local\". Host must be an IP address",
		},
		{
			name:        "invalid scheme",
			url:         "ftp://192.168.0.100",
			wantErr:     true,
			wantErrText: "invalid URL scheme: \"ftp\". Scheme must be http or https",
		},
		{
			name:        "empty url",
			url:         "",
			wantErr:     true,
			wantErrText: "invalid URL scheme: \"\". Scheme must be http or https",
		},
		{
			name:        "missing scheme",
			url:         "192.168.0.100",
			wantErr:     true,
			wantErrText: "invalid URL scheme: \"\". Scheme must be http or https",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateWiiMURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateWiiMURL(%q) err = %v, wantErr %v", tt.url, err != nil, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.wantErrText {
				t.Errorf("validateWiiMURL(%q) errText = %q, wantErrText %q", tt.url, err.Error(), tt.wantErrText)
			}
		})
	}
}
