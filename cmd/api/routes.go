package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	if app.spotifyClient != nil {
		mux.HandleFunc("GET /spotify/login", app.spotifyLogin)
		mux.HandleFunc("GET /spotify/callback", app.spotifyCallback)
	}

	if app.lastfmClient != nil {
		mux.HandleFunc("GET /lastfm/login", app.lastFmLogin)
		mux.HandleFunc("GET /lastfm/callback", app.lastFmCallback)
	}

	return logRequest(mux)
}
