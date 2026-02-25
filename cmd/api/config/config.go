package config

import (
	"github.com/degeens/scrobblet/internal/clients"
	"github.com/degeens/scrobblet/internal/clients/csv"
	"github.com/degeens/scrobblet/internal/clients/lastfm"
	"github.com/degeens/scrobblet/internal/clients/listenbrainz"
	"github.com/degeens/scrobblet/internal/clients/spotify"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

func LoadConfig() (*Config, error) {
	port := getEnv(envPort, defaultPort)
	dataPath := getEnv(envDataPath, defaultDataPath)
	logLevel := getEnv(envLogLevel, defaultLogLevel)

	rateLimitRate, err := getIntEnv(envRateLimitRate, defaultRateLimitRate)
	if err != nil {
		return nil, err
	}

	rateLimitBurst, err := getIntEnv(envRateLimitBurst, defaultRateLimitBurst)
	if err != nil {
		return nil, err
	}

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

	return &Config{
		Port:           port,
		DataPath:       dataPath,
		LogLevel:       logLevel,
		RateLimitRate:  rateLimitRate,
		RateLimitBurst: rateLimitBurst,
		Source:         sourceType,
		Targets:        targetTypes,
		Clients:        clientsConfig,
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

	err = validateRedirectURL(redirectURL, "/spotify/callback")
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

	err = validateRedirectURL(redirectURL, "/lastfm/callback")
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
