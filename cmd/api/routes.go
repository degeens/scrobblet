package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	if app.spotifyClient != nil {
		mux.HandleFunc("GET /spotify/login", app.login)
		mux.HandleFunc("GET /spotify/callback", app.callback)
	}

	if app.lastfmClient != nil {
		mux.HandleFunc("GET /lastfm/login", app.lastfmLogin)
		mux.HandleFunc("GET /lastfm/callback", app.lastfmCallback)
	}

	return logRequest(mux)
}
