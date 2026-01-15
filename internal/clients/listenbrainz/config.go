package listenbrainz

type Config struct {
	URL   string // Optional: custom base URL (defaults to https://api.listenbrainz.org)
	Token string
}
