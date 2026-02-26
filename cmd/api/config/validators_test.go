package config

import (
	"reflect"
	"testing"

	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

func TestError(t *testing.T) {
	t.Errorf("This test is failing")
}

func TestValidateSource(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    sources.SourceType
		wantErr bool
	}{
		{
			name:    "valid lowercase",
			source:  "spotify",
			want:    sources.SourceSpotify,
			wantErr: false,
		},
		{
			name:    "valid titlecase",
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
			got, err := validateSource(tt.source)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateSource(%q) err = %v, wantErr %v", tt.source, err != nil, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("validateSource(%q) = %q, want %q", tt.source, got, tt.want)
			}
		})
	}
}

func TestValidateTargets(t *testing.T) {
	tests := []struct {
		name    string
		targets string
		want    []targets.TargetType
		wantErr bool
	}{
		{
			name:    "single valid lowercase",
			targets: "koito",
			want:    []targets.TargetType{targets.TargetKoito},
			wantErr: false,
		},
		{
			name:    "multiple valid lowercase",
			targets: "koito,csv",
			want:    []targets.TargetType{targets.TargetKoito, targets.TargetCSV},
			wantErr: false,
		},
		{
			name:    "multiple valid titlecase",
			targets: "Koito,CSV",
			want:    []targets.TargetType{targets.TargetKoito, targets.TargetCSV},
			wantErr: false,
		},
		{
			name:    "multiple valid with surrounding spaces",
			targets: " Koito , CSV ",
			want:    []targets.TargetType{targets.TargetKoito, targets.TargetCSV},
			wantErr: false,
		},
		{
			name:    "unknown returns error",
			targets: "unknown",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "unknown among valid returns error",
			targets: "Koito,CSV,unknown",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "duplicate returns error",
			targets: "Koito,Koito",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty string returns error",
			targets: "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateTargets(tt.targets)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateTargets(%q) err = %v, wantErr %v", tt.targets, err != nil, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateTargets(%q) = %v, want %v", tt.targets, got, tt.want)
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
