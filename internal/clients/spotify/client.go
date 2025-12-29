package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const (
	baseURL  = "https://api.spotify.com/v1"
	authURL  = "https://accounts.spotify.com/authorize"
	tokenURL = "https://accounts.spotify.com/api/token"
)

type Client struct {
	baseURL      string
	httpClient   *http.Client
	oauth2Config *oauth2.Config
	oauth2Token  *oauth2.Token
}

func NewClient(clientID, clientSecret, redirectURL string) *Client {
	c := &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
		oauth2Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"user-read-currently-playing"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  authURL,
				TokenURL: tokenURL,
			},
		},
	}

	c.loadToken()

	return c
}

func (c *Client) GetCurrentlyPlayingTrack() (*CurrentlyPlayingTrack, error) {
	url := c.baseURL + "/me/player/currently-playing"

	if err := c.refreshTokenIfNeeded(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.oauth2Token.AccessToken)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		var currentlyPlaying CurrentlyPlayingTrack
		if err := json.NewDecoder(res.Body).Decode(&currentlyPlaying); err != nil {
			return nil, err
		}
		return &currentlyPlaying, nil
	case http.StatusNoContent:
		return nil, nil
	default:
		return nil, fmt.Errorf("Failed to get currently playing track from Spotify (HTTP status code %d)", res.StatusCode)
	}
}
