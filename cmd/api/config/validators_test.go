package config

import (
	"reflect"
	"testing"

	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

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
				t.Fatalf("validateSource(%q) error = %v, wantErr %v", tt.source, err, tt.wantErr)
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
				t.Fatalf("validateTargets(%q) error = %v, wantErr %v", tt.targets, err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateTargets(%q) = %v, want %v", tt.targets, got, tt.want)
			}
		})
	}
}
