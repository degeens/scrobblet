package main

import (
	"fmt"
	"os"
)

type config struct {
	port    string
	spotify spotifyConfig
	koito   koitoConfig
}

type spotifyConfig struct {
	clientID     string
	clientSecret string
	redirectURL  string
}

type koitoConfig struct {
	url   string
	token string
}

const (
	// Default values
	defaultPort = "8080"

	// Environment variable keys
	envPort                = "SCROBBLET_PORT"
	envSpotifyClientID     = "SPOTIFY_CLIENT_ID"
	envSpotifyClientSecret = "SPOTIFY_CLIENT_SECRET"
	envSpotifyRedirectURL  = "SPOTIFY_REDIRECT_URL"
	envKoitoURL            = "KOITO_URL"
	envKoitoToken          = "KOITO_TOKEN"
)

func loadConfig() (*config, error) {
	port := getEnv(envPort, defaultPort)

	spotify, err := loadSpotifyConfig()
	if err != nil {
		return nil, err
	}

	koito, err := loadKoitoConfig()
	if err != nil {
		return nil, err
	}

	return &config{
		port:    port,
		spotify: spotify,
		koito:   koito,
	}, nil
}

func loadKoitoConfig() (koitoConfig, error) {
	url, err := getRequiredEnv(envKoitoURL)
	if err != nil {
		return koitoConfig{}, err
	}

	token, err := getRequiredEnv(envKoitoToken)
	if err != nil {
		return koitoConfig{}, err
	}

	return koitoConfig{
		url:   url,
		token: token,
	}, nil
}

func loadSpotifyConfig() (spotifyConfig, error) {
	clientID, err := getRequiredEnv(envSpotifyClientID)
	if err != nil {
		return spotifyConfig{}, err
	}

	clientSecret, err := getRequiredEnv(envSpotifyClientSecret)
	if err != nil {
		return spotifyConfig{}, err
	}

	redirectURL, err := getRequiredEnv(envSpotifyRedirectURL)
	if err != nil {
		return spotifyConfig{}, err
	}

	return spotifyConfig{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
	}, nil
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func getRequiredEnv(key string) (string, error) {
	value := os.Getenv(key)

	if value == "" {
		return "", fmt.Errorf("Required environment variable %s is not set", key)
	}

	return value, nil
}
