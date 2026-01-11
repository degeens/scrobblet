package clients

import (
	"github.com/degeens/scrobblet/internal/clients/koito"
	"github.com/degeens/scrobblet/internal/clients/lastfm"
	"github.com/degeens/scrobblet/internal/clients/listenbrainz"
	"github.com/degeens/scrobblet/internal/clients/spotify"
)

type Config struct {
	Spotify      spotify.Config
	Koito        koito.Config
	ListenBrainz listenbrainz.Config
	LastFm       lastfm.Config
}
