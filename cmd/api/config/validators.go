package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

func validateSource(source string) (sources.SourceType, error) {
	switch source {
	case string(sources.SourceSpotify):
		return sources.SourceSpotify, nil
	default:
		return "", fmt.Errorf("invalid source: %s. Valid sources are: %s", source, sources.SourceSpotify)
	}
}

func validateTargets(targetsString string) ([]targets.TargetType, error) {
	targetStrings := strings.Split(targetsString, ",")
	targetTypes := make([]targets.TargetType, 0, len(targetStrings))
	seen := make(map[targets.TargetType]bool)

	for _, targetString := range targetStrings {
		targetString = strings.TrimSpace(targetString)

		targetType, err := validateTarget(targetString)
		if err != nil {
			return nil, err
		}

		if seen[targetType] {
			return nil, fmt.Errorf("duplicate target type: %s. Multiple targets of the same type are not supported", targetType)
		}
		seen[targetType] = true

		targetTypes = append(targetTypes, targetType)
	}

	if len(targetTypes) == 0 {
		return nil, fmt.Errorf("no targets specified in %s", envTargets)
	}

	return targetTypes, nil
}

func validateTarget(target string) (targets.TargetType, error) {
	switch target {
	case string(targets.TargetKoito):
		return targets.TargetKoito, nil
	case string(targets.TargetMaloja):
		return targets.TargetMaloja, nil
	case string(targets.TargetListenBrainz):
		return targets.TargetListenBrainz, nil
	case string(targets.TargetLastFm):
		return targets.TargetLastFm, nil
	case string(targets.TargetCSV):
		return targets.TargetCSV, nil
	default:
		return "", fmt.Errorf("invalid target: %s. Valid targets are: %s, %s, %s, %s, %s", target, targets.TargetKoito, targets.TargetMaloja, targets.TargetListenBrainz, targets.TargetLastFm, targets.TargetCSV)
	}
}

func validateRedirectURL(pathPrefix, redirectURL string) error {
	parsedURL, err := url.Parse(redirectURL)
	if err != nil {
		return fmt.Errorf("invalid redirect URL: %w", err)
	}

	if parsedURL.Path != pathPrefix+"/callback" {
		return fmt.Errorf("invalid redirect URL path: %s. Path must be /callback", parsedURL.Path)
	}

	return nil
}
