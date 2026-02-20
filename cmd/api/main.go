package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/degeens/scrobblet/cmd/api/config"
	"github.com/degeens/scrobblet/internal/clients/lastfm"
	"github.com/degeens/scrobblet/internal/clients/spotify"
	"github.com/degeens/scrobblet/internal/scrobbler"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

var version = "undefined" // Will be overridden at build time

type application struct {
	source         *sources.Source
	targets        *[]targets.Target
	spotifyClient  *spotify.Client
	lastfmClient   *lastfm.Client
	authStateStore *authStateStore
	config         *config.Config
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	setDefaultLogger(cfg.LogLevel)

	slog.Info("Starting Scrobblet", "version", version)

	sourceClient, source, err := sources.New(cfg.Source, cfg.Clients)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	targetClients, targets, err := targets.NewMultiple(cfg.Targets, cfg.Clients, version)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		source:         &source,
		targets:        &targets,
		authStateStore: newAuthStateStore(),
		config:         cfg,
	}
	if spotifyClient, ok := sourceClient.(*spotify.Client); ok {
		app.spotifyClient = spotifyClient
	}
	for _, client := range targetClients {
		if lastfmClient, ok := client.(*lastfm.Client); ok {
			app.lastfmClient = lastfmClient
		}
	}

	scrobbler := scrobbler.NewScrobbler(source, targets)

	go scrobbler.Start()
	slog.Info("Scrobbler started")

	slog.Info("Listening on port :" + cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, app.routes())
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func setDefaultLogger(logLevel string) {
	var level slog.Level

	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)
}
