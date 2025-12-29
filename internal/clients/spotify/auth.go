package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
)

func (c *Client) GetAuthCodeURL() string {
	return c.oauth2Config.AuthCodeURL("")
}

func (c *Client) ExchangeAuthCodeForToken(ctx context.Context, code string) error {
	token, err := c.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return err
	}

	c.oauth2Token = token

	return c.saveToken()
}

func (c *Client) refreshTokenIfNeeded() error {
	if c.oauth2Token == nil {
		return errors.New("Not authenticated, log in via /login")
	}

	// Check if token is expired or about to expire (within 5 minutes)
	if time.Now().Add(5 * time.Minute).Before(c.oauth2Token.Expiry) {
		return nil // Token is still valid
	}

	// Use oauth2 TokenSource to automatically refresh
	token, err := c.oauth2Config.TokenSource(context.Background(), c.oauth2Token).Token()
	if err != nil {
		return err
	}

	c.oauth2Token = token

	return c.saveToken()
}

func (c *Client) loadToken() error {
	tokenFilePath, err := getTokenFilePath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(tokenFilePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}

		return err
	}

	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return err
	}

	c.oauth2Token = &token

	return nil
}

func (c *Client) saveToken() error {
	data, err := json.Marshal(c.oauth2Token)
	if err != nil {
		return err
	}

	tokenFilePath, err := getTokenFilePath()
	if err != nil {
		return err
	}

	tokenDir := filepath.Dir(tokenFilePath)
	if err := os.MkdirAll(tokenDir, 0700); err != nil {
		return err
	}

	return os.WriteFile(tokenFilePath, data, 0600)
}

func getTokenFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "scrobblet", "spotify_token.json"), nil
}
