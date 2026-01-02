package lastfm

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL = "http://ws.audioscrobbler.com/2.0/"
)

type Client struct {
	baseURL        string
	httpClient     *http.Client
	apiKey         string
	secret         string
	sessionKey     string
	sessionKeyPath string
}

func NewClient(apiKey, secret, sessionKey, dataPath string) *Client {
	client := &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiKey:     apiKey,
		secret:     secret,
		sessionKey: sessionKey,
	}

	if dataPath != "" {
		client.sessionKeyPath = dataPath + "/lastfm_session.json"
		// Try to load existing session key if none provided
		if sessionKey == "" {
			client.loadSessionKey()
		}
	}

	return client
}

// generateSignature creates an MD5 signature for API calls
// Parameters must be sorted alphabetically by key name
func (c *Client) generateSignature(params map[string]string) string {
	// Get keys and sort them
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build the signature string
	var sigString strings.Builder
	for _, k := range keys {
		sigString.WriteString(k)
		sigString.WriteString(params[k])
	}
	sigString.WriteString(c.secret)

	// Generate MD5 hash
	hash := md5.Sum([]byte(sigString.String()))
	return fmt.Sprintf("%x", hash)
}

// UpdateNowPlaying sends a now playing notification to Last.fm
func (c *Client) UpdateNowPlaying(req *UpdateNowPlayingRequest) error {
	params := map[string]string{
		"method":  "track.updateNowPlaying",
		"artist":  req.Artist,
		"track":   req.Track,
		"api_key": c.apiKey,
		"sk":      c.sessionKey,
		"format":  "json",
	}

	if req.Album != "" {
		params["album"] = req.Album
	}
	if req.Duration > 0 {
		params["duration"] = strconv.Itoa(req.Duration)
	}
	if req.TrackNumber > 0 {
		params["trackNumber"] = strconv.Itoa(req.TrackNumber)
	}
	if req.AlbumArtist != "" {
		params["albumArtist"] = req.AlbumArtist
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
		return fmt.Errorf("failed to send now playing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var nowPlayingResp NowPlayingResponse
	if err := json.Unmarshal(body, &nowPlayingResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	if nowPlayingResp.Error != nil {
		msg := "Unknown error"
		if nowPlayingResp.Message != nil {
			msg = *nowPlayingResp.Message
		}
		return fmt.Errorf("Last.fm API error (code %d): %s", *nowPlayingResp.Error, msg)
	}

	// Check if track was ignored
	if nowPlayingResp.NowPlaying.IgnoredMessage.Code != "" && nowPlayingResp.NowPlaying.IgnoredMessage.Code != "0" {
		return fmt.Errorf("track was ignored (code %s): %s",
			nowPlayingResp.NowPlaying.IgnoredMessage.Code,
			nowPlayingResp.NowPlaying.IgnoredMessage.Text)
	}

	return nil
}

// Scrobble sends one or more scrobbles to Last.fm
func (c *Client) Scrobble(scrobbles []ScrobbleRequest) error {
	if len(scrobbles) == 0 {
		return fmt.Errorf("no scrobbles to submit")
	}
	if len(scrobbles) > 50 {
		return fmt.Errorf("cannot submit more than 50 scrobbles at once")
	}

	params := map[string]string{
		"method":  "track.scrobble",
		"api_key": c.apiKey,
		"sk":      c.sessionKey,
		"format":  "json",
	}

	// Add scrobble parameters with array notation
	for i, scrobble := range scrobbles {
		idx := ""
		if len(scrobbles) > 1 {
			idx = fmt.Sprintf("[%d]", i)
		}

		params["artist"+idx] = scrobble.Artist
		params["track"+idx] = scrobble.Track
		params["timestamp"+idx] = strconv.FormatInt(scrobble.Timestamp, 10)

		if scrobble.Album != "" {
			params["album"+idx] = scrobble.Album
		}
		if scrobble.Duration > 0 {
			params["duration"+idx] = strconv.Itoa(scrobble.Duration)
		}
		if scrobble.TrackNumber > 0 {
			params["trackNumber"+idx] = strconv.Itoa(scrobble.TrackNumber)
		}
		if scrobble.AlbumArtist != "" {
			params["albumArtist"+idx] = scrobble.AlbumArtist
		}
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
		return fmt.Errorf("failed to send scrobble request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var scrobbleResp ScrobbleResponse
	if err := json.Unmarshal(body, &scrobbleResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	if scrobbleResp.Error != nil {
		msg := "Unknown error"
		if scrobbleResp.Message != nil {
			msg = *scrobbleResp.Message
		}
		return fmt.Errorf("Last.fm API error (code %d): %s", *scrobbleResp.Error, msg)
	}

	// Check if any scrobbles were ignored
	ignored := 0
	if scrobbleResp.Scrobbles.Attr.Ignored != "" {
		ignored, _ = strconv.Atoi(scrobbleResp.Scrobbles.Attr.Ignored)
	}
	if ignored > 0 {
		var ignoredReasons []string
		for _, scrobble := range scrobbleResp.Scrobbles.Scrobble {
			if scrobble.IgnoredMessage.Code != "" && scrobble.IgnoredMessage.Code != "0" {
				ignoredReasons = append(ignoredReasons, fmt.Sprintf("'%s - %s' (code %s): %s",
					scrobble.Artist.Text, scrobble.Track.Text, scrobble.IgnoredMessage.Code, scrobble.IgnoredMessage.Text))
			}
		}
		if len(ignoredReasons) > 0 {
			return fmt.Errorf("%d scrobble(s) ignored: %s", ignored, strings.Join(ignoredReasons, "; "))
		}
	}

	return nil
}
