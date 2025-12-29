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
