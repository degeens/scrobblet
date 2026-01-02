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
	authURL = "https://www.last.fm/api/auth"
)

type session struct {
	Key string `json:"session_key"`
}

type getSessionResponse struct {
	Session struct {
		Key string `json:"key"`
	} `json:"session"`
	Error   *int    `json:"error"`
	Message *string `json:"message"`
}

func (c *Client) GetAuthURL() (string, error) {
	params := url.Values{}
	params.Set("api_key", c.apiKey)
	params.Set("cb", c.redirectURL)

	url, err := url.Parse(authURL)
	if err != nil {
		return "", err
	}

	url.RawQuery = params.Encode()

	return url.String(), nil
}

func (c *Client) ExchangeTokenForSession(ctx context.Context, token string) error {
	params := map[string]string{
		"method":  "auth.getSession",
		"api_key": c.apiKey,
		"token":   token,
		"format":  "json",
	}

	params["api_sig"] = c.generateSignature(params)

	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}

	res, err := c.httpClient.PostForm(c.baseURL, formData)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var getSessionResp getSessionResponse
	if err := json.Unmarshal(body, &getSessionResp); err != nil {
		return err
	}

	if getSessionResp.Error != nil {
		return fmt.Errorf("Failed to get session from Last.fm. API error (code %d): %s", *getSessionResp.Error, *getSessionResp.Message)
	}

	c.session = &session{Key: getSessionResp.Session.Key}

	return c.saveSessionKey()
}

func (c *Client) loadSessionKey() error {
	data, err := os.ReadFile(c.sessionPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}

		return err
	}

	var session session
	if err := json.Unmarshal(data, &session); err != nil {
		return err
	}

	c.session = &session

	return nil
}

func (c *Client) saveSessionKey() error {
	data, err := json.Marshal(c.session)
	if err != nil {
		return err
	}

	sessionDir := filepath.Dir(c.sessionPath)
	if err := os.MkdirAll(sessionDir, 0700); err != nil {
		return err
	}

	return os.WriteFile(c.sessionPath, data, 0600)
}
