package main

import (
	"context"
	"net/http"
)

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	url := app.spotifyClient.GetAuthCodeURL()

	http.Redirect(w, r, url, http.StatusFound)
}

func (app *application) callback(w http.ResponseWriter, r *http.Request) {
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
	w.Write([]byte("Authentication successful!"))
}

func (app *application) lastfmLogin(w http.ResponseWriter, r *http.Request) {
	// Get the callback URL from the config (passed via environment variable)
	callbackURL := r.URL.Query().Get("callback")
	if callbackURL == "" {
		// Use default from request
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		callbackURL = scheme + "://" + r.Host + "/lastfm/callback"
	}

	url := app.lastfmClient.GetAuthURL(callbackURL)

	http.Redirect(w, r, url, http.StatusFound)
}

func (app *application) lastfmCallback(w http.ResponseWriter, r *http.Request) {
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
	w.Write([]byte("Last.fm authentication successful! You can now start scrobbling."))
}
