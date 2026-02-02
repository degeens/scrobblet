package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/degeens/scrobblet/internal/clients"
	"github.com/degeens/scrobblet/internal/clients/csv"
	"github.com/degeens/scrobblet/internal/clients/lastfm"
	"github.com/degeens/scrobblet/internal/clients/listenbrainz"
	"github.com/degeens/scrobblet/internal/clients/spotify"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

type config struct {
	port     string
	dataPath string
	logLevel string
	source   sources.SourceType
	targets  []targets.TargetType
	clients  clients.Config
}

const (
	// Default values
	defaultPort     = "7276"
	defaultDataPath = "/etc/scrobblet"
	defaultLogLevel = "INFO"

	// Environment variable keys
	envPort                = "SCROBBLET_PORT"
	envDataPath            = "SCROBBLET_DATA_PATH"
	envLogLevel            = "SCROBBLET_LOG_LEVEL"
	envSource              = "SCROBBLET_SOURCE"
	envTargets             = "SCROBBLET_TARGETS"
	envSpotifyClientID     = "SPOTIFY_CLIENT_ID"
	envSpotifyClientSecret = "SPOTIFY_CLIENT_SECRET"
	envSpotifyRedirectURL  = "SPOTIFY_REDIRECT_URL"
	envKoitoURL            = "KOITO_URL"
	envKoitoToken          = "KOITO_TOKEN"
	envMalojaURL           = "MALOJA_URL"
	envMalojaToken         = "MALOJA_TOKEN"
	envListenBrainzURL     = "LISTENBRAINZ_URL"
	envListenBrainzToken   = "LISTENBRAINZ_TOKEN"
	envLastFmAPIKey        = "LASTFM_API_KEY"
	envLastFmSharedSecret  = "LASTFM_SHARED_SECRET"
	envLastFmRedirectURL   = "LASTFM_REDIRECT_URL"
	envCSVFilePath         = "CSV_FILE_PATH"
)

func loadConfig() (*config, error) {
	port := getEnv(envPort, defaultPort)
	dataPath := getEnv(envDataPath, defaultDataPath)
	logLevel := getEnv(envLogLevel, defaultLogLevel)

	source, err := getRequiredEnv(envSource)
	if err != nil {
		return nil, err
	}

	sourceType, err := validateSource(source)
	if err != nil {
		return nil, err
	}

	targets, err := getRequiredEnv(envTargets)
	if err != nil {
		return nil, err
	}

	targetTypes, err := validateTargets(targets)
	if err != nil {
		return nil, err
	}

	clientsConfig, err := loadClientsConfig(sourceType, targetTypes, dataPath)
	if err != nil {
		return nil, err
	}

	return &config{
		port:     port,
		dataPath: dataPath,
		logLevel: logLevel,
		source:   sourceType,
		targets:  targetTypes,
		clients:  clientsConfig,
	}, nil
}

func loadClientsConfig(sourceType sources.SourceType, targetTypes []targets.TargetType, dataPath string) (clients.Config, error) {
	var spotifyConfig spotify.Config
	var koitoConfig listenbrainz.Config
	var malojaConfig listenbrainz.Config
	var listenBrainzConfig listenbrainz.Config
	var lastfmConfig lastfm.Config
	var csvConfig csv.Config
	var err error

	if sourceType == sources.SourceSpotify {
		spotifyConfig, err = loadSpotifyConfig(dataPath)
		if err != nil {
			return clients.Config{}, err
		}
	}

	// Load configs for each target type
	for _, targetType := range targetTypes {
		switch targetType {
		case targets.TargetKoito:
			koitoConfig, err = loadListenBrainzConfig(targetType)
			if err != nil {
				return clients.Config{}, err
			}
		case targets.TargetMaloja:
			malojaConfig, err = loadListenBrainzConfig(targetType)
			if err != nil {
				return clients.Config{}, err
			}
		case targets.TargetListenBrainz:
			listenBrainzConfig, err = loadListenBrainzConfig(targetType)
			if err != nil {
				return clients.Config{}, err
			}
		case targets.TargetLastFm:
			lastfmConfig, err = loadLastFmConfig(dataPath)
			if err != nil {
				return clients.Config{}, err
			}
		case targets.TargetCSV:
			csvConfig, err = loadCSVConfig(dataPath)
			if err != nil {
				return clients.Config{}, err
			}
		}
	}

	return clients.Config{
		Spotify:      spotifyConfig,
		Koito:        koitoConfig,
		Maloja:       malojaConfig,
		ListenBrainz: listenBrainzConfig,
		LastFm:       lastfmConfig,
		CSV:          csvConfig,
	}, nil
}

func loadSpotifyConfig(dataPath string) (spotify.Config, error) {
	clientID, err := getRequiredEnv(envSpotifyClientID)
	if err != nil {
		return spotify.Config{}, err
	}

	clientSecret, err := getRequiredEnv(envSpotifyClientSecret)
	if err != nil {
		return spotify.Config{}, err
	}

	redirectURL, err := getRequiredEnv(envSpotifyRedirectURL)
	if err != nil {
		return spotify.Config{}, err
	}

	err = validateRedirectURL("/spotify", redirectURL)
	if err != nil {
		return spotify.Config{}, err
	}

	return spotify.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		DataPath:     dataPath,
	}, nil
}

func loadListenBrainzConfig(targetType targets.TargetType) (listenbrainz.Config, error) {
	var token, baseURL string
	var err error

	switch targetType {
	case targets.TargetKoito:
		token, err = getRequiredEnv(envKoitoToken)
		if err != nil {
			return listenbrainz.Config{}, err
		}
		baseURL, err = getRequiredEnv(envKoitoURL)
		if err != nil {
			return listenbrainz.Config{}, err
		}
		baseURL, err = setURLPath(baseURL, "/apis/listenbrainz")
		if err != nil {
			return listenbrainz.Config{}, err
		}
	case targets.TargetMaloja:
		token, err = getRequiredEnv(envMalojaToken)
		if err != nil {
			return listenbrainz.Config{}, err
		}
		baseURL, err = getRequiredEnv(envMalojaURL)
		if err != nil {
			return listenbrainz.Config{}, err
		}
		baseURL, err = setURLPath(baseURL, "/apis/listenbrainz")
		if err != nil {
			return listenbrainz.Config{}, err
		}
	default:
		token, err = getRequiredEnv(envListenBrainzToken)
		if err != nil {
			return listenbrainz.Config{}, err
		}
		baseURL = getEnv(envListenBrainzURL, "")
	}

	return listenbrainz.Config{
		URL:   baseURL,
		Token: token,
	}, nil
}

func loadLastFmConfig(dataPath string) (lastfm.Config, error) {
	apiKey, err := getRequiredEnv(envLastFmAPIKey)
	if err != nil {
		return lastfm.Config{}, err
	}

	sharedSecret, err := getRequiredEnv(envLastFmSharedSecret)
	if err != nil {
		return lastfm.Config{}, err
	}

	redirectURL, err := getRequiredEnv(envLastFmRedirectURL)
	if err != nil {
		return lastfm.Config{}, err
	}

	err = validateRedirectURL("/lastfm", redirectURL)
	if err != nil {
		return lastfm.Config{}, err
	}

	return lastfm.Config{
		APIKey:       apiKey,
		SharedSecret: sharedSecret,
		RedirectURL:  redirectURL,
		DataPath:     dataPath,
	}, nil
}

func loadCSVConfig(dataPath string) (csv.Config, error) {
	filePath := getEnv(envCSVFilePath, dataPath+"/scrobbles.csv")

	return csv.Config{
		FilePath: filePath,
	}, nil
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func getRequiredEnv(key string) (string, error) {
	value := os.Getenv(key)

	if value == "" {
		return "", fmt.Errorf("required environment variable %s is not set", key)
	}

	return value, nil
}

func validateSource(source string) (sources.SourceType, error) {
	switch source {
	case string(sources.SourceSpotify):
		return sources.SourceSpotify, nil
	default:
		return "", fmt.Errorf("invalid source: %s. Valid sources are: %s", source, sources.SourceSpotify)
	}
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

func setURLPath(rawURL, path string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// Clear any existing path and set to the provided path
	parsedURL.Path = path
	parsedURL.RawPath = ""

	return parsedURL.String(), nil
}
