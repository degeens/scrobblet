package koito

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

func NewClient(baseURL string, token string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
		token:      token,
	}
}

func (c *Client) SubmitListens(request *SubmitListens) error {
	url := c.baseURL + "/apis/listenbrainz/1/submit-listens"

	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Token "+c.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("failed to submit listens to Koito (HTTP %d): %s", res.StatusCode, string(body))
}
