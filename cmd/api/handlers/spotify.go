package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/degeens/scrobblet/cmd/api/utils"
	"github.com/degeens/scrobblet/internal/clients/spotify"
)

func SpotifyLogin(spotifyClient *spotify.Client, authStateStore *utils.AuthStateStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, err := authStateStore.Generate("spotify")
		if err != nil {
			slog.Error("Failed to generate state", "error", err)
			http.Error(w, authenticationFailed, http.StatusInternalServerError)
			return
		}

		url := spotifyClient.GetAuthCodeURL(state)

		http.Redirect(w, r, url, http.StatusFound)
	}
}

func SpotifyCallback(spotifyClient *spotify.Client, authStateStore *utils.AuthStateStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			slog.Error("Failed to parse form", "error", err)
			http.Error(w, authenticationFailed, http.StatusInternalServerError)
			return
		}

		state := r.Form.Get("state")
		if err := authStateStore.Validate("spotify", state); err != nil {
			http.Error(w, authenticationFailed, http.StatusBadRequest)
			return
		}

		code := r.Form.Get("code")
		if code == "" {
			http.Error(w, authenticationFailed, http.StatusBadRequest)
			return
		}

		err = spotifyClient.ExchangeAuthCodeForToken(context.Background(), code)
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
}
