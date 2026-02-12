package main

import (
	"context"
	"log/slog"
	"net/http"
)

const (
	authenticationSucceeded = "Authentication successful! Feel free to close this browser tab."
	authenticationFailed    = "Authentication failed. Please try again."
)

func (app *application) spotifyLogin(w http.ResponseWriter, r *http.Request) {
	state, err := app.authStateStore.Generate("spotify")
	if err != nil {
		slog.Error("Failed to generate state", "error", err)
		http.Error(w, authenticationFailed, http.StatusInternalServerError)
		return
	}

	url := app.spotifyClient.GetAuthCodeURL(state)

	http.Redirect(w, r, url, http.StatusFound)
}

func (app *application) spotifyCallback(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		slog.Error("Failed to parse form", "error", err)
		http.Error(w, authenticationFailed, http.StatusInternalServerError)
		return
	}

	state := r.Form.Get("state")
	if err := app.authStateStore.Validate("spotify", state); err != nil {
		http.Error(w, authenticationFailed, http.StatusBadRequest)
		return
	}

	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, authenticationFailed, http.StatusBadRequest)
		return
	}

	err = app.spotifyClient.ExchangeAuthCodeForToken(context.Background(), code)
	if err != nil {
		slog.Error("Failed to exchange code for token", "error", err)
		http.Error(w, authenticationFailed, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(authenticationSucceeded))
	if err != nil {
		slog.Error("Failed to write response", "error", err)
	}
}

func (app *application) lastFmLogin(w http.ResponseWriter, r *http.Request) {
	url, err := app.lastfmClient.GetAuthURL()
	if err != nil {
		slog.Error("Failed to get auth URL", "error", err)
		http.Error(w, authenticationFailed, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (app *application) lastFmCallback(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		slog.Error("Failed to parse form", "error", err)
		http.Error(w, authenticationFailed, http.StatusInternalServerError)
		return
	}

	token := r.Form.Get("token")
	if token == "" {
		http.Error(w, authenticationFailed, http.StatusBadRequest)
		return
	}

	err = app.lastfmClient.ExchangeTokenForSession(context.Background(), token)
	if err != nil {
		slog.Error("Failed to exchange token for session", "error", err)
		http.Error(w, authenticationFailed, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(authenticationSucceeded))
	if err != nil {
		slog.Error("Failed to write response", "error", err)
	}
}
