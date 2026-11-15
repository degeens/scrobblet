package api

import (
	"fmt"
	"net/http"

	"github.com/degeens/scrobblet/internal/api/handlers"
	"github.com/degeens/scrobblet/internal/api/middleware"
	"github.com/degeens/scrobblet/internal/api/utils"
	"github.com/degeens/scrobblet/internal/clients/lastfm"
	"github.com/degeens/scrobblet/internal/clients/spotify"
	"github.com/degeens/scrobblet/internal/config"
	"github.com/degeens/scrobblet/internal/metrics"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewHandler(source sources.Source, targetInstances []targets.Instance, sourceClient any, cfg *config.Config, metrics *metrics.Metrics) http.Handler {
	apiMux := http.NewServeMux()
	authStateStore := utils.NewAuthStateStore()

	spotifyClient := getSpotifyClient(sourceClient)
	if spotifyClient != nil {
		apiMux.HandleFunc("GET /spotify/login", handlers.SpotifyLogin(spotifyClient, authStateStore))
		apiMux.HandleFunc("GET /spotify/callback", handlers.SpotifyCallback(spotifyClient, authStateStore))
	}

	for _, instance := range targetInstances {
		lastfmClient := getLastFmClient(instance)
		if lastfmClient == nil {
			continue
		}

		if instance.Name == lastfm.DefaultInstanceName {
			apiMux.HandleFunc("GET /lastfm/login", handlers.LastFmLogin(lastfmClient))
			apiMux.HandleFunc("GET /lastfm/callback", handlers.LastFmCallback(lastfmClient))
		} else {
			apiMux.HandleFunc(fmt.Sprintf("GET /lastfm/%s/login", instance.Name), handlers.LastFmLogin(lastfmClient))
			apiMux.HandleFunc(fmt.Sprintf("GET /lastfm/%s/callback", instance.Name), handlers.LastFmCallback(lastfmClient))
		}
	}

	targetList := make([]targets.Target, 0, len(targetInstances))
	for _, instance := range targetInstances {
		targetList = append(targetList, instance.Target)
	}

	rootMux := http.NewServeMux()
	rootMux.Handle("GET /health", handlers.Health(source, targetList))
	rootMux.Handle("GET /metrics", promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{}))
	rootMux.Handle("/api/", http.StripPrefix("/api", middleware.LogRequest(middleware.RateLimit(cfg.RateLimitRate, cfg.RateLimitBurst)(apiMux))))

	return rootMux
}

func getSpotifyClient(sourceClient any) *spotify.Client {
	spotifyClient, _ := sourceClient.(*spotify.Client)

	return spotifyClient
}

func getLastFmClient(instance targets.Instance) *lastfm.Client {
	if instance.Type != targets.TargetLastFm {
		return nil
	}

	lastfmClient, _ := instance.Client.(*lastfm.Client)

	return lastfmClient
}
