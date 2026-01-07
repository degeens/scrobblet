package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

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

	// TokenSource automatically handles token refresh when needed
	token, err := c.oauth2Config.TokenSource(context.Background(), c.oauth2Token).Token()
	if err != nil {
		return err
	}

	// Only save if the token was actually refreshed
	if token.AccessToken != c.oauth2Token.AccessToken {
		c.oauth2Token = token
		return c.saveToken()
	}

	return nil
}

func (c *Client) loadToken() error {
	data, err := os.ReadFile(c.oauth2TokenPath)
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

	tokenDir := filepath.Dir(c.oauth2TokenPath)
	if err := os.MkdirAll(tokenDir, 0700); err != nil {
		return err
	}

	return os.WriteFile(c.oauth2TokenPath, data, 0600)
}
