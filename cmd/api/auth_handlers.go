package main

import (
	"context"
	"net/http"
)

func (app *application) spotifyLogin(w http.ResponseWriter, r *http.Request) {
	url := app.spotifyClient.GetAuthCodeURL()

	http.Redirect(w, r, url, http.StatusFound)
}

func (app *application) spotifyCallback(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}

	err = app.spotifyClient.ExchangeAuthCodeForToken(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Spotify authentication successful! Feel free to close this browser tab."))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *application) lastFmLogin(w http.ResponseWriter, r *http.Request) {
	url, err := app.lastfmClient.GetAuthURL()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (app *application) lastFmCallback(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := r.Form.Get("token")
	if token == "" {
		http.Error(w, "Authorization token not found", http.StatusBadRequest)
		return
	}

	err = app.lastfmClient.ExchangeTokenForSession(context.Background(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Last.fm authentication successful! Feel free to close this browser tab."))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
