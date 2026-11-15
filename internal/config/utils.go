package config

import (
	"fmt"
	"net/url"
)

func setURLPath(rawURL, path string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// Clear any existing path and set to the provided path
	parsedURL.Path = path
	parsedURL.RawPath = ""

	return parsedURL.String(), nil
}
