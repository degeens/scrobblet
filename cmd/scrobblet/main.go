package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/degeens/scrobblet/internal/api"
	"github.com/degeens/scrobblet/internal/config"
	"github.com/degeens/scrobblet/internal/metrics"
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

	sourceClient, source, err := sources.New(cfg.Source)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	targetInstances, err := targets.NewMultiple(cfg.Targets, version)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	targetList := make([]targets.Target, 0, len(targetInstances))
	for _, instance := range targetInstances {
		targetList = append(targetList, instance.Target)
	}

	metrics := metrics.New()
	metrics.BuildInfo.WithLabelValues(version).Set(1)

	scrobbler := scrobbler.NewScrobbler(source, targetList, metrics)

	go scrobbler.Start()
	slog.Info("Scrobbler started")

	slog.Info("Listening on port :" + cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, api.NewHandler(source, targetInstances, sourceClient, cfg, metrics))
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
