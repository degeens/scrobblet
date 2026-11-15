package config

import (
	"fmt"
	"net/netip"
	"net/url"
	"regexp"
	"strings"

	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

// Target names are used in URL paths (Last.fm auth routes) and file names
// (session and CSV files), so they are restricted to a safe charset.
var targetNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

func validateSourceType(source string) (sources.SourceType, error) {
	source = strings.TrimSpace(source)

	switch strings.ToLower(source) {
	case strings.ToLower(string(sources.SourceSpotify)):
		return sources.SourceSpotify, nil
	case strings.ToLower(string(sources.SourceWiiM)):
		return sources.SourceWiiM, nil
	default:
		return "", fmt.Errorf("invalid source type: %s. Valid source types are: %s, %s", source, sources.SourceSpotify, sources.SourceWiiM)
	}
}

func validateTargetType(target string) (targets.TargetType, error) {
	switch strings.ToLower(strings.TrimSpace(target)) {
	case strings.ToLower(string(targets.TargetKoito)):
		return targets.TargetKoito, nil
	case strings.ToLower(string(targets.TargetMaloja)):
		return targets.TargetMaloja, nil
	case strings.ToLower(string(targets.TargetListenBrainz)):
		return targets.TargetListenBrainz, nil
	case strings.ToLower(string(targets.TargetLastFm)):
		return targets.TargetLastFm, nil
	case strings.ToLower(string(targets.TargetCSV)):
		return targets.TargetCSV, nil
	default:
		return "", fmt.Errorf("invalid target type: %s. Valid target types are: %s, %s, %s, %s, %s", target, targets.TargetKoito, targets.TargetMaloja, targets.TargetListenBrainz, targets.TargetLastFm, targets.TargetCSV)
	}
}

func validateTargetName(name string) error {
	if !targetNamePattern.MatchString(name) {
		return fmt.Errorf("invalid name: %q. Name must start with a lowercase letter or digit and may only contain lowercase letters, digits and dashes", name)
	}

	return nil
}

func validateRedirectURL(redirectURL, validPath string) error {
	parsedURL, err := url.Parse(redirectURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("invalid URL scheme: %q. Scheme must be http or https", parsedURL.Scheme)
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("invalid URL: host must not be empty")
	}

	if parsedURL.Path != validPath {
		return fmt.Errorf("invalid URL path: %q. Path must be %q", parsedURL.Path, validPath)
	}

	return nil
}

// The WiiM client skips TLS certificate verification, so it must only ever
// connect to devices on the local network.
func validateWiiMURL(rawURL string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("invalid URL scheme: %q. Scheme must be http or https", parsedURL.Scheme)
	}

	host := parsedURL.Hostname()
	addr, err := netip.ParseAddr(host)
	if err != nil {
		return fmt.Errorf("invalid URL host: %q. Host must be an IP address", host)
	}

	if !addr.IsPrivate() {
		return fmt.Errorf("invalid URL host: %q. Host must be a local network IP address", host)
	}

	return nil
}
