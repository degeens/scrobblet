package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
	"gopkg.in/yaml.v3"
)

// Config is the validated, construction-ready application configuration.
type Config struct {
	Port           string
	DataPath       string
	LogLevel       string
	RateLimitRate  int
	RateLimitBurst int
	Source         sources.Spec
	Targets        []targets.Spec
}

// configFile mirrors the YAML configuration file schema.
type configFile struct {
	Port      *int             `yaml:"port"`
	DataPath  string           `yaml:"dataPath"`
	LogLevel  string           `yaml:"logLevel"`
	RateLimit *rateLimitConfig `yaml:"rateLimit"`
	Source    *sourceEntry     `yaml:"source"`
	Targets   []targetEntry    `yaml:"targets"`
}

type rateLimitConfig struct {
	Rate  *int `yaml:"rate"`
	Burst *int `yaml:"burst"`
}

type sourceHeader struct {
	Type string `yaml:"type"`
}

type spotifySourceOptions struct {
	sourceHeader `yaml:",inline"`
	ClientID     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	RedirectURL  string `yaml:"redirectUrl"`
}

type wiimSourceOptions struct {
	sourceHeader `yaml:",inline"`
	URL          string `yaml:"url"`
}

// sourceEntry dispatches on the source type into typed per-type options.
type sourceEntry struct {
	Type    string
	Spotify *spotifySourceOptions
	WiiM    *wiimSourceOptions
}

func (e *sourceEntry) UnmarshalYAML(node *yaml.Node) error {
	var header sourceHeader
	if err := node.Decode(&header); err != nil {
		return err
	}

	e.Type = header.Type

	switch strings.ToLower(strings.TrimSpace(header.Type)) {
	case strings.ToLower(string(sources.SourceSpotify)):
		var opts spotifySourceOptions
		if err := strictDecode(node, &opts); err != nil {
			return err
		}
		e.Spotify = &opts
	case strings.ToLower(string(sources.SourceWiiM)):
		var opts wiimSourceOptions
		if err := strictDecode(node, &opts); err != nil {
			return err
		}
		e.WiiM = &opts
	default:
		// Unknown source types are reported during validation
	}

	return nil
}

type targetHeader struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
}

type listenBrainzTargetOptions struct {
	targetHeader `yaml:",inline"`
	URL          string `yaml:"url"`
	Token        string `yaml:"token"`
}

type lastFmTargetOptions struct {
	targetHeader `yaml:",inline"`
	APIKey       string `yaml:"apiKey"`
	SharedSecret string `yaml:"sharedSecret"`
	RedirectURL  string `yaml:"redirectUrl"`
}

type csvTargetOptions struct {
	targetHeader `yaml:",inline"`
	FilePath     string `yaml:"filePath"`
}

// targetEntry dispatches on the target type into typed per-type options.
type targetEntry struct {
	Type         string
	Name         string
	ListenBrainz *listenBrainzTargetOptions
	LastFm       *lastFmTargetOptions
	CSV          *csvTargetOptions
}

func (e *targetEntry) UnmarshalYAML(node *yaml.Node) error {
	var header targetHeader
	if err := node.Decode(&header); err != nil {
		return err
	}

	e.Type = header.Type
	e.Name = header.Name

	switch strings.ToLower(strings.TrimSpace(header.Type)) {
	case strings.ToLower(string(targets.TargetListenBrainz)),
		strings.ToLower(string(targets.TargetMaloja)),
		strings.ToLower(string(targets.TargetKoito)):
		var opts listenBrainzTargetOptions
		if err := strictDecode(node, &opts); err != nil {
			return err
		}
		e.ListenBrainz = &opts
	case strings.ToLower(string(targets.TargetLastFm)):
		var opts lastFmTargetOptions
		if err := strictDecode(node, &opts); err != nil {
			return err
		}
		e.LastFm = &opts
	case strings.ToLower(string(targets.TargetCSV)):
		var opts csvTargetOptions
		if err := strictDecode(node, &opts); err != nil {
			return err
		}
		e.CSV = &opts
	default:
		// Unknown target types are reported during validation
	}

	return nil
}

// strictDecode decodes a YAML node into out, rejecting unknown fields.
// node.Decode alone does not support KnownFields, so the node is re-encoded
// and decoded with a strict decoder.
func strictDecode(node *yaml.Node, out any) error {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	if err := encoder.Encode(node); err != nil {
		return fmt.Errorf("failed to re-encode YAML node: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to re-encode YAML node: %w", err)
	}

	decoder := yaml.NewDecoder(&buf)
	decoder.KnownFields(true)

	return decoder.Decode(out)
}
