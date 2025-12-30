package clients

import (
	"github.com/degeens/scrobblet/internal/clients/koito"
	"github.com/degeens/scrobblet/internal/clients/spotify"
)

type Config struct {
	Spotify spotify.Config
	Koito   koito.Config
}
