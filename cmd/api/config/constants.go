package config

// Default values
const (
	defaultPort           = "7276"
	defaultDataPath       = "/etc/scrobblet"
	defaultLogLevel       = "INFO"
	defaultRateLimitRate  = 10
	defaultRateLimitBurst = 100
)

// Environment variable keys
const (
	envPort                = "SCROBBLET_PORT"
	envDataPath            = "SCROBBLET_DATA_PATH"
	envLogLevel            = "SCROBBLET_LOG_LEVEL"
	envRateLimitRate       = "SCROBBLET_RATE_LIMIT_RATE"
	envRateLimitBurst      = "SCROBBLET_RATE_LIMIT_BURST"
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
