package config

// Default values
const (
	defaultPort           = 7276
	defaultDataPath       = "/etc/scrobblet"
	defaultLogLevel       = "INFO"
	defaultRateLimitRate  = 10
	defaultRateLimitBurst = 100

	defaultConfigFileName    = "scrobblet.yaml"
	defaultListenBrainzURL   = "https://api.listenbrainz.org"
	listenBrainzAPIPath      = "/apis/listenbrainz"
	spotifyCallbackPath      = "/api/spotify/callback"
	lastFmCallbackPath       = "/api/lastfm/callback"
	lastFmNamedCallbackTmpl  = "/api/lastfm/%s/callback"
	defaultCSVFileName       = "scrobbles.csv"
	namedCSVFileNameTemplate = "scrobbles-%s.csv"
)

// Environment variable keys
const (
	envConfig = "SCROBBLET_CONFIG"
)
