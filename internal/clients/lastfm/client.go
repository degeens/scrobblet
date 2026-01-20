package lastfm

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL = "https://ws.audioscrobbler.com/2.0/"
)

type Client struct {
	baseURL      string
	userAgent    string
	httpClient   *http.Client
	apiKey       string
	sharedSecret string
	redirectURL  string
	session      *session
	sessionPath  string
}

func NewClient(apiKey, sharedSecret, redirectURL, dataPath, scrobbletVersion string) *Client {
	c := &Client{
		baseURL:      baseURL,
		userAgent:    fmt.Sprintf("Scrobblet/%s", scrobbletVersion),
		httpClient:   &http.Client{Timeout: 15 * time.Second},
		apiKey:       apiKey,
		sharedSecret: sharedSecret,
		redirectURL:  redirectURL,
	}

	c.sessionPath = filepath.Join(dataPath, "lastfm_session.json")

	c.loadSessionKey()

	return c
}

func (c *Client) UpdateNowPlaying(request *UpdateNowPlayingRequest) error {
	if c.session == nil {
		return errors.New("not authenticated with Last.fm, log in via /lastfm/login")
	}

	params := map[string]string{
		"method":      "track.updateNowPlaying",
		"artist":      request.Artist,
		"track":       request.Track,
		"album":       request.Album,
		"trackNumber": strconv.Itoa(request.TrackNumber),
		"duration":    strconv.Itoa(request.Duration),
		"api_key":     c.apiKey,
		"sk":          c.session.Key,
		"format":      "json",
	}

	params["api_sig"] = c.generateSignature(params)

	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}

	req, err := http.NewRequest("POST", c.baseURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var updateNowPlayingResp UpdateNowPlayingResponse
	if err := json.Unmarshal(body, &updateNowPlayingResp); err != nil {
		return err
	}

	if updateNowPlayingResp.Error != nil {
		return fmt.Errorf("failed to update now playing track on Last.fm. API error (code %d): %s", *updateNowPlayingResp.Error, *updateNowPlayingResp.Message)
	}

	if updateNowPlayingResp.NowPlaying.IgnoredMessage.Code != "0" {
		return fmt.Errorf("failed to update now playing track on Last.fm. Track ignored (code %s): %s", updateNowPlayingResp.NowPlaying.IgnoredMessage.Code, updateNowPlayingResp.NowPlaying.IgnoredMessage.Text)
	}

	return nil
}

func (c *Client) Scrobble(requests []ScrobbleRequest) error {
	if c.session == nil {
		return errors.New("not authenticated with Last.fm, log in via /lastfm/login")
	}

	if len(requests) == 0 {
		return fmt.Errorf("no scrobbles to submit")
	}
	if len(requests) > 50 {
		return fmt.Errorf("cannot submit more than 50 scrobbles per batch")
	}

	params := map[string]string{
		"method":  "track.scrobble",
		"api_key": c.apiKey,
		"sk":      c.session.Key,
		"format":  "json",
	}

	for i, request := range requests {
		idx := ""
		if len(requests) > 1 {
			idx = fmt.Sprintf("[%d]", i)
		}

		params["artist"+idx] = request.Artist
		params["track"+idx] = request.Track
		params["timestamp"+idx] = strconv.FormatInt(request.Timestamp, 10)
		params["album"+idx] = request.Album
		params["trackNumber"+idx] = strconv.Itoa(request.TrackNumber)
		params["duration"+idx] = strconv.Itoa(request.Duration)
	}

	params["api_sig"] = c.generateSignature(params)

	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}

	req, err := http.NewRequest("POST", c.baseURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var scrobbleResp ScrobbleResponse
	if err := json.Unmarshal(body, &scrobbleResp); err != nil {
		return err
	}

	if scrobbleResp.Error != nil {
		return fmt.Errorf("failed to submit scrobble to Last.fm. API error (code %d): %s", *scrobbleResp.Error, *scrobbleResp.Message)
	}

	if scrobbleResp.Scrobbles.Attr.Ignored > 0 {
		var ignoredReasons []string
		for _, scrobble := range scrobbleResp.Scrobbles.Scrobble {
			if scrobble.IgnoredMessage.Code != "0" {
				ignoredReasons = append(ignoredReasons, fmt.Sprintf("%s - %s (code %s): %s", scrobble.Artist.Text, scrobble.Track.Text, scrobble.IgnoredMessage.Code, scrobble.IgnoredMessage.Text))
			}
		}

		if len(ignoredReasons) > 0 {
			return fmt.Errorf("failed to submit scrobble to Last.fm. %d track(s) ignored: %s", scrobbleResp.Scrobbles.Attr.Ignored, strings.Join(ignoredReasons, ", "))
		}
	}

	return nil
}

func (c *Client) generateSignature(params map[string]string) string {
	// Params to exclude from signature
	excludeParams := map[string]bool{
		"format": true,
	}

	// Exclude and sort params
	keys := make([]string, 0, len(params))
	for k := range params {
		if !excludeParams[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// Build the signature
	var signatureBuilder strings.Builder
	for _, k := range keys {
		signatureBuilder.WriteString(k)
		signatureBuilder.WriteString(params[k])
	}
	signatureBuilder.WriteString(c.sharedSecret)

	// Generate an MD5 hash of the signature
	hash := md5.Sum([]byte(signatureBuilder.String()))

	return fmt.Sprintf("%x", hash)
}
