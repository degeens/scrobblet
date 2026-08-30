package api

import (
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

func NewHandler(source sources.Source, targets []targets.Target, sourceClient any, targetClients []any, cfg *config.Config, metrics *metrics.Metrics) http.Handler {
	apiMux := http.NewServeMux()
	authStateStore := utils.NewAuthStateStore()

	spotifyClient := getSpotifyClient(sourceClient)
	if spotifyClient != nil {
		apiMux.HandleFunc("GET /spotify/login", handlers.SpotifyLogin(spotifyClient, authStateStore))
		apiMux.HandleFunc("GET /spotify/callback", handlers.SpotifyCallback(spotifyClient, authStateStore))
	}

	lastfmClient := getLastFmClient(targetClients)
	if lastfmClient != nil {
		apiMux.HandleFunc("GET /lastfm/login", handlers.LastFmLogin(lastfmClient))
		apiMux.HandleFunc("GET /lastfm/callback", handlers.LastFmCallback(lastfmClient))
	}

	rootMux := http.NewServeMux()
	rootMux.Handle("GET /health", handlers.Health(source, targets))
	rootMux.Handle("GET /metrics", promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{}))
	rootMux.Handle("/api/", http.StripPrefix("/api", middleware.LogRequest(middleware.RateLimit(cfg.RateLimitRate, cfg.RateLimitBurst)(apiMux))))

	return rootMux
}

func getSpotifyClient(sourceClient any) *spotify.Client {
	spotifyClient, _ := sourceClient.(*spotify.Client)

	return spotifyClient
}

func getLastFmClient(targetClients []any) *lastfm.Client {
	var lastfmClient *lastfm.Client
	for _, client := range targetClients {
		if c, ok := client.(*lastfm.Client); ok {
			lastfmClient = c
			break
		}
	}

	return lastfmClient
}
