package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/degeens/scrobblet/cmd/api/config"
	"github.com/degeens/scrobblet/cmd/api/utils"
	"github.com/degeens/scrobblet/internal/scrobbler"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

var version = "undefined" // Will be overridden at build time

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

	scrobbler := scrobbler.NewScrobbler(source, targets)

	authStateStore := utils.NewAuthStateStore()

	go scrobbler.Start()
	slog.Info("Scrobbler started")

	slog.Info("Listening on port :" + cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, routes(source, targets, sourceClient, targetClients, cfg, authStateStore))
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
