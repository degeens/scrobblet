package main

import (
	"net/http"

	"github.com/degeens/scrobblet/cmd/api/middleware"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", app.health)

	if app.spotifyClient != nil {
		mux.HandleFunc("GET /spotify/login", app.spotifyLogin)
		mux.HandleFunc("GET /spotify/callback", app.spotifyCallback)
	}

	if app.lastfmClient != nil {
		mux.HandleFunc("GET /lastfm/login", app.lastFmLogin)
		mux.HandleFunc("GET /lastfm/callback", app.lastFmCallback)
	}

	rate := app.config.RateLimitRate
	burst := app.config.RateLimitBurst

	return middleware.RateLimit(rate, burst)(middleware.LogRequest(mux))
}
