package main

import (
	"log/slog"
	"net/http"
	"os"

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

	sourceClient, source, err := sources.New(cfg.source, cfg.clients)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	_, target, err := targets.New(cfg.target, cfg.clients)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	app := &application{}
	if spotifyClient, ok := sourceClient.(*spotify.Client); ok {
		app.spotifyClient = spotifyClient
	}

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
