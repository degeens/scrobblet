package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/degeens/scrobblet/internal/clients"
	"github.com/degeens/scrobblet/internal/clients/koito"
	"github.com/degeens/scrobblet/internal/clients/listenbrainz"
	"github.com/degeens/scrobblet/internal/clients/spotify"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

type config struct {
	port     string
	dataPath string
	source   sources.SourceType
	target   targets.TargetType
	clients  clients.Config
}

const (
	// Default values
	defaultPort     = "7276"
	defaultDataPath = "/etc/scrobblet"

	// Environment variable keys
	envPort                = "SCROBBLET_PORT"
	envDataPath            = "SCROBBLET_DATA_PATH"
	envSource              = "SCROBBLET_SOURCE"
	envTarget              = "SCROBBLET_TARGET"
	envSpotifyClientID     = "SPOTIFY_CLIENT_ID"
	envSpotifyClientSecret = "SPOTIFY_CLIENT_SECRET"
	envSpotifyRedirectURL  = "SPOTIFY_REDIRECT_URL"
	envKoitoURL            = "KOITO_URL"
	envKoitoToken          = "KOITO_TOKEN"
	envListenBrainzToken   = "LISTENBRAINZ_TOKEN"
)

func loadConfig() (*config, error) {
	port := getEnv(envPort, defaultPort)
	dataPath := getEnv(envDataPath, defaultDataPath)

	source, err := getRequiredEnv(envSource)
	if err != nil {
		return nil, err
	}

	sourceType, err := validateSource(source)
	if err != nil {
		return nil, err
	}

	target, err := getRequiredEnv(envTarget)
	if err != nil {
		return nil, err
	}

	targetType, err := validateTarget(target)
	if err != nil {
		return nil, err
	}

	clientsConfig, err := loadClientsConfig(sourceType, targetType, dataPath)
	if err != nil {
		return nil, err
	}

	return &config{
		port:     port,
		dataPath: dataPath,
		source:   sourceType,
		target:   targetType,
		clients:  clientsConfig,
	}, nil
}

func loadClientsConfig(sourceType sources.SourceType, targetType targets.TargetType, dataPath string) (clients.Config, error) {
	var spotifyConfig spotify.Config
	var koitoConfig koito.Config
	var listenbrainzConfig listenbrainz.Config
	var err error

	if sourceType == sources.SourceSpotify {
		spotifyConfig, err = loadSpotifyConfig(dataPath)
		if err != nil {
			return clients.Config{}, err
		}
	}

	if targetType == targets.TargetKoito {
		koitoConfig, err = loadKoitoConfig()
		if err != nil {
			return clients.Config{}, err
		}
	}

	if targetType == targets.TargetListenBrainz {
		listenbrainzConfig, err = loadListenBrainzConfig()
		if err != nil {
			return clients.Config{}, err
		}
	}

	return clients.Config{
		Spotify:      spotifyConfig,
		Koito:        koitoConfig,
		ListenBrainz: listenbrainzConfig,
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

	err = validateRedirectURL(redirectURL)
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

func loadKoitoConfig() (koito.Config, error) {
	url, err := getRequiredEnv(envKoitoURL)
	if err != nil {
		return koito.Config{}, err
	}

	token, err := getRequiredEnv(envKoitoToken)
	if err != nil {
		return koito.Config{}, err
	}

	return koito.Config{
		URL:   url,
		Token: token,
	}, nil
}

func loadListenBrainzConfig() (listenbrainz.Config, error) {
	token, err := getRequiredEnv(envListenBrainzToken)
	if err != nil {
		return listenbrainz.Config{}, err
	}

	return listenbrainz.Config{
		Token: token,
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
		return "", fmt.Errorf("Required environment variable %s is not set", key)
	}

	return value, nil
}

func validateSource(source string) (sources.SourceType, error) {
	switch source {
	case string(sources.SourceSpotify):
		return sources.SourceSpotify, nil
	default:
		return "", fmt.Errorf("Invalid source: %s. Valid sources are: %s.", source, sources.SourceSpotify)
	}
}

func validateTarget(target string) (targets.TargetType, error) {
	switch target {
	case string(targets.TargetKoito):
		return targets.TargetKoito, nil
	case string(targets.TargetListenBrainz):
		return targets.TargetListenBrainz, nil
	default:
		return "", fmt.Errorf("Invalid target: %s. Valid targets are: %s, %s.", target, targets.TargetKoito, targets.TargetListenBrainz)
	}
}

func validateRedirectURL(redirectURL string) error {
	parsedURL, err := url.Parse(redirectURL)
	if err != nil {
		return fmt.Errorf("Invalid redirect URL: %w", err)
	}

	if parsedURL.Path != "/callback" {
		return fmt.Errorf("Invalid redirect URL path: %s. Path must be /callback.", parsedURL.Path)
	}

	return nil
}
