package main

import (
	"net/http"

	"github.com/degeens/scrobblet/cmd/api/config"
	"github.com/degeens/scrobblet/cmd/api/handlers"
	"github.com/degeens/scrobblet/cmd/api/middleware"
	"github.com/degeens/scrobblet/cmd/api/utils"
	"github.com/degeens/scrobblet/internal/clients/lastfm"
	"github.com/degeens/scrobblet/internal/clients/spotify"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

func routes(source sources.Source, targets []targets.Target, sourceClient any, targetClients []any, config *config.Config, authStateStore *utils.AuthStateStore) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.Health(source, targets))

	if spotifyClient, ok := sourceClient.(*spotify.Client); ok {
		mux.HandleFunc("GET /spotify/login", handlers.SpotifyLogin(spotifyClient, authStateStore))
		mux.HandleFunc("GET /spotify/callback", handlers.SpotifyCallback(spotifyClient, authStateStore))
	}

	for _, client := range targetClients {
		if lastfmClient, ok := client.(*lastfm.Client); ok {
			mux.HandleFunc("GET /lastfm/login", handlers.LastFmLogin(lastfmClient))
			mux.HandleFunc("GET /lastfm/callback", handlers.LastFmCallback(lastfmClient))
			break
		}
	}

	rate := config.RateLimitRate
	burst := config.RateLimitBurst

	return middleware.RateLimit(rate, burst)(middleware.LogRequest(mux))
}
