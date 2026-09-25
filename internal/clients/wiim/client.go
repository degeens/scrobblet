package wiim

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	c := &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 5 * time.Second, Transport: transport},
	}

	return c
}

func (c *Client) GetPlayerStatus() (*PlayerStatus, error) {
	url := c.baseURL + "/httpapi.asp?command=getPlayerStatus"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode == http.StatusOK {
		var playerStatus PlayerStatus
		if err := json.Unmarshal(body, &playerStatus); err != nil {
			return nil, err
		}
		return &playerStatus, nil
	} else {
		return nil, fmt.Errorf("failed to get currently playing track from WiiM (HTTP %d): %s", res.StatusCode, string(body))
	}
}
