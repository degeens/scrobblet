package clients

import (
	"github.com/degeens/scrobblet/internal/clients/csv"
	"github.com/degeens/scrobblet/internal/clients/lastfm"
	"github.com/degeens/scrobblet/internal/clients/listenbrainz"
	"github.com/degeens/scrobblet/internal/clients/spotify"
)

type Config struct {
	Spotify      spotify.Config
	ListenBrainz listenbrainz.Config
	LastFm       lastfm.Config
	CSV          csv.Config
}
