package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/degeens/scrobblet/internal/clients/koito"
	"github.com/degeens/scrobblet/internal/clients/spotify"
	"github.com/degeens/scrobblet/internal/scrobbler"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

type application struct {
	spotifyClient *spotify.Client
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("Config loaded")

	spotifyClient := spotify.NewClient(cfg.spotify.clientID, cfg.spotify.clientSecret, cfg.spotify.redirectURL)
	koitoClient := koito.NewClient(cfg.koito.url, cfg.koito.token)

	app := &application{
		spotifyClient: spotifyClient,
	}

	source := sources.NewSpotifySource(spotifyClient)
	target := targets.NewKoitoTarget(koitoClient)
	scrobbler := scrobbler.NewScrobbler(source, target)

	go scrobbler.Start()
	slog.Info("Scrobbler started")

	slog.Info("Listening on port :" + cfg.port)
	err = http.ListenAndServe(":"+cfg.port, app.routes())
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
