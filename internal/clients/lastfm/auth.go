package lastfm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
)

const (
	authURL = "http://www.last.fm/api/auth"
)

// GetAuthURL returns the URL for user authorization
func (c *Client) GetAuthURL(callbackURL string) string {
	params := url.Values{}
	params.Set("api_key", c.apiKey)
	if callbackURL != "" {
		params.Set("cb", callbackURL)
	}
	return authURL + "?" + params.Encode()
}

// GetSession exchanges an authorization token for a session key
func (c *Client) GetSession(ctx context.Context, token string) (string, error) {
	params := map[string]string{
		"method":  "auth.getSession",
		"api_key": c.apiKey,
		"token":   token,
		"format":  "json",
	}

	// Generate signature
	params["api_sig"] = c.generateSignature(params)

	// Make the request
	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}

	resp, err := c.httpClient.PostForm(c.baseURL, formData)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var sessionResp struct {
		Session struct {
			Name       string `json:"name"`
			Key        string `json:"key"`
			Subscriber int    `json:"subscriber"`
		} `json:"session"`
		Error   *int    `json:"error"`
		Message *string `json:"message"`
	}

	if err := json.Unmarshal(body, &sessionResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	if sessionResp.Error != nil {
		msg := "Unknown error"
		if sessionResp.Message != nil {
			msg = *sessionResp.Message
		}
		return "", fmt.Errorf("Last.fm API error (code %d): %s", *sessionResp.Error, msg)
	}

	if sessionResp.Session.Key == "" {
		return "", fmt.Errorf("session key not found in response")
	}

	return sessionResp.Session.Key, nil
}

// ExchangeTokenForSession exchanges an authorization token for a session key and stores it
func (c *Client) ExchangeTokenForSession(ctx context.Context, token string) error {
	sessionKey, err := c.GetSession(ctx, token)
	if err != nil {
		return err
	}

	c.sessionKey = sessionKey

	return c.saveSessionKey()
}

// loadSessionKey loads the session key from disk
func (c *Client) loadSessionKey() error {
	if c.sessionKeyPath == "" {
		return nil // No path configured, skip loading
	}

	data, err := os.ReadFile(c.sessionKeyPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}

	var session struct {
		SessionKey string `json:"session_key"`
	}

	if err := json.Unmarshal(data, &session); err != nil {
		return err
	}

	c.sessionKey = session.SessionKey

	return nil
}

// saveSessionKey saves the session key to disk
func (c *Client) saveSessionKey() error {
	if c.sessionKeyPath == "" {
		return nil // No path configured, skip saving
	}

	session := struct {
		SessionKey string `json:"session_key"`
	}{
		SessionKey: c.sessionKey,
	}

	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	sessionDir := filepath.Dir(c.sessionKeyPath)
	if err := os.MkdirAll(sessionDir, 0700); err != nil {
		return err
	}

	return os.WriteFile(c.sessionKeyPath, data, 0600)
}
