package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	if app.spotifyClient != nil {
		mux.HandleFunc("GET /login", app.login)
		mux.HandleFunc("GET /callback", app.callback)
	}

	return logRequest(mux)
}
