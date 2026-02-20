package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func getRequiredEnv(key string) (string, error) {
	value := os.Getenv(key)

	if value == "" {
		return "", fmt.Errorf("required environment variable %s is not set", key)
	}

	return value, nil
}

func getIntEnv(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue, nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer value for %s", key)
	}

	return intValue, nil
}

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
