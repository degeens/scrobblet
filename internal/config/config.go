package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/degeens/scrobblet/internal/clients/csv"
	"github.com/degeens/scrobblet/internal/clients/lastfm"
	"github.com/degeens/scrobblet/internal/clients/listenbrainz"
	"github.com/degeens/scrobblet/internal/clients/spotify"
	"github.com/degeens/scrobblet/internal/clients/wiim"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
	"gopkg.in/yaml.v3"
)

func LoadConfig() (*Config, error) {
	configPath := os.Getenv(envConfig)
	if configPath == "" {
		configPath = filepath.Join(defaultDataPath, defaultConfigFileName)
	}

	raw, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %w", configPath, err)
	}

	// Expand ${ENV_VAR} references once, before parsing, so tokens can stay in
	// environment variables or Docker secrets.
	expanded := os.ExpandEnv(string(raw))

	var file configFile
	decoder := yaml.NewDecoder(strings.NewReader(expanded))
	decoder.KnownFields(true)
	if err := decoder.Decode(&file); err != nil {
		return nil, fmt.Errorf("failed to parse config file %q: %w", configPath, err)
	}

	return toConfig(&file)
}

func toConfig(file *configFile) (*Config, error) {
	port := defaultPort
	if file.Port != nil {
		port = *file.Port
	}

	dataPath := file.DataPath
	if dataPath == "" {
		dataPath = defaultDataPath
	}

	logLevel := file.LogLevel
	if logLevel == "" {
		logLevel = defaultLogLevel
	}

	rateLimitRate := defaultRateLimitRate
	rateLimitBurst := defaultRateLimitBurst
	if file.RateLimit != nil {
		if file.RateLimit.Rate != nil {
			rateLimitRate = *file.RateLimit.Rate
		}
		if file.RateLimit.Burst != nil {
			rateLimitBurst = *file.RateLimit.Burst
		}
	}

	source, err := toSourceSpec(file.Source, dataPath)
	if err != nil {
		return nil, err
	}

	targetSpecs, err := toTargetSpecs(file.Targets, dataPath)
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:           strconv.Itoa(port),
		DataPath:       dataPath,
		LogLevel:       logLevel,
		RateLimitRate:  rateLimitRate,
		RateLimitBurst: rateLimitBurst,
		Source:         source,
		Targets:        targetSpecs,
	}, nil
}

func toSourceSpec(entry *sourceEntry, dataPath string) (sources.Spec, error) {
	if entry == nil {
		return sources.Spec{}, fmt.Errorf("source is required")
	}

	sourceType, err := validateSourceType(entry.Type)
	if err != nil {
		return sources.Spec{}, err
	}

	switch sourceType {
	case sources.SourceSpotify:
		opts := entry.Spotify
		if opts.ClientID == "" {
			return sources.Spec{}, fmt.Errorf("source (Spotify): clientId is required")
		}
		if opts.ClientSecret == "" {
			return sources.Spec{}, fmt.Errorf("source (Spotify): clientSecret is required")
		}
		if opts.RedirectURL == "" {
			return sources.Spec{}, fmt.Errorf("source (Spotify): redirectUrl is required")
		}
		if err := validateRedirectURL(opts.RedirectURL, spotifyCallbackPath); err != nil {
			return sources.Spec{}, fmt.Errorf("source (Spotify): redirectUrl: %w", err)
		}

		return sources.Spec{
			Type: sourceType,
			Spotify: &spotify.Config{
				ClientID:     opts.ClientID,
				ClientSecret: opts.ClientSecret,
				RedirectURL:  opts.RedirectURL,
				DataPath:     dataPath,
			},
		}, nil
	case sources.SourceWiiM:
		opts := entry.WiiM
		if opts.URL == "" {
			return sources.Spec{}, fmt.Errorf("source (WiiM): url is required")
		}
		if err := validateWiiMURL(opts.URL); err != nil {
			return sources.Spec{}, fmt.Errorf("source (WiiM): url: %w", err)
		}

		return sources.Spec{
			Type: sourceType,
			WiiM: &wiim.Config{BaseURL: opts.URL},
		}, nil
	default:
		return sources.Spec{}, fmt.Errorf("unknown source type: %s", sourceType)
	}
}

func toTargetSpecs(entries []targetEntry, dataPath string) ([]targets.Spec, error) {
	if len(entries) == 0 {
		return nil, fmt.Errorf("at least one target is required")
	}

	specs := make([]targets.Spec, 0, len(entries))
	seen := make(map[string]int)
	csvFilePaths := make(map[string]int)

	for i, entry := range entries {
		spec, csvFilePath, err := toTargetSpec(entry, i, dataPath)
		if err != nil {
			return nil, err
		}

		key := string(spec.Type) + "/" + spec.Name
		if j, ok := seen[key]; ok {
			return nil, fmt.Errorf("targets[%d] (%s %q): duplicate target, already defined at targets[%d]. Targets of the same type must have a unique name", i, spec.Type, spec.Name, j)
		}
		seen[key] = i

		if csvFilePath != "" {
			if j, ok := csvFilePaths[csvFilePath]; ok {
				return nil, fmt.Errorf("targets[%d] (%s %q): filePath %q conflicts with targets[%d]. CSV targets must write to different files", i, spec.Type, spec.Name, csvFilePath, j)
			}
			csvFilePaths[csvFilePath] = i
		}

		specs = append(specs, spec)
	}

	return specs, nil
}

//gocyclo:ignore
func toTargetSpec(entry targetEntry, index int, dataPath string) (targets.Spec, string, error) {
	targetType, err := validateTargetType(entry.Type)
	if err != nil {
		return targets.Spec{}, "", fmt.Errorf("targets[%d]: %w", index, err)
	}

	name := entry.Name
	if name == "" {
		name = defaultTargetName(targetType)
	}

	if err := validateTargetName(name); err != nil {
		return targets.Spec{}, "", fmt.Errorf("targets[%d] (%s): %w", index, targetType, err)
	}

	describe := func() string {
		return fmt.Sprintf("targets[%d] (%s %q)", index, targetType, name)
	}

	spec := targets.Spec{Type: targetType, Name: name}
	csvFilePath := ""

	switch targetType {
	case targets.TargetListenBrainz, targets.TargetMaloja, targets.TargetKoito:
		opts := entry.ListenBrainz
		if opts.Token == "" {
			return targets.Spec{}, "", fmt.Errorf("%s: token is required", describe())
		}

		baseURL := opts.URL
		switch targetType {
		case targets.TargetListenBrainz:
			if baseURL == "" {
				baseURL = defaultListenBrainzURL
			}
		case targets.TargetMaloja, targets.TargetKoito:
			// Koito and Maloja are ListenBrainz-compatible services with a fixed API path
			if baseURL == "" {
				return targets.Spec{}, "", fmt.Errorf("%s: url is required", describe())
			}
			baseURL, err = setURLPath(baseURL, listenBrainzAPIPath)
			if err != nil {
				return targets.Spec{}, "", fmt.Errorf("%s: url: %w", describe(), err)
			}
		}

		spec.ListenBrainz = &listenbrainz.Config{URL: baseURL, Token: opts.Token}
	case targets.TargetLastFm:
		opts := entry.LastFm
		if opts.APIKey == "" {
			return targets.Spec{}, "", fmt.Errorf("%s: apiKey is required", describe())
		}
		if opts.SharedSecret == "" {
			return targets.Spec{}, "", fmt.Errorf("%s: sharedSecret is required", describe())
		}
		if opts.RedirectURL == "" {
			return targets.Spec{}, "", fmt.Errorf("%s: redirectUrl is required", describe())
		}

		callbackPath := lastFmCallbackPath
		if name != lastfm.DefaultInstanceName {
			callbackPath = fmt.Sprintf(lastFmNamedCallbackTmpl, name)
		}
		if err := validateRedirectURL(opts.RedirectURL, callbackPath); err != nil {
			return targets.Spec{}, "", fmt.Errorf("%s: redirectUrl: %w", describe(), err)
		}

		spec.LastFm = &lastfm.Config{
			APIKey:       opts.APIKey,
			SharedSecret: opts.SharedSecret,
			RedirectURL:  opts.RedirectURL,
			DataPath:     dataPath,
		}
	case targets.TargetCSV:
		opts := entry.CSV
		filePath := opts.FilePath
		if filePath == "" {
			fileName := defaultCSVFileName
			if name != defaultTargetName(targets.TargetCSV) {
				fileName = fmt.Sprintf(namedCSVFileNameTemplate, name)
			}
			filePath = filepath.Join(dataPath, fileName)
		}

		csvFilePath = filePath
		spec.CSV = &csv.Config{FilePath: filePath}
	}

	return spec, csvFilePath, nil
}

func defaultTargetName(targetType targets.TargetType) string {
	return strings.ToLower(string(targetType))
}
