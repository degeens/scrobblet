package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

var (
	ErrInvalidState = errors.New("invalid or expired OAuth state")
)

type authState struct {
	value     string
	expiresAt time.Time
}

type AuthStateStore struct {
	mu     sync.Mutex
	states map[string]authState // key is service name (e.g., "spotify", "lastfm")
}

func NewAuthStateStore() *AuthStateStore {
	return &AuthStateStore{
		states: make(map[string]authState),
	}
}

// Generates a new random state parameter for a specific service and stores it.
// Stores one state parameter per service; generating a new state overwrites the previous one.
func (m *AuthStateStore) Generate(service string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(b)

	m.mu.Lock()
	defer m.mu.Unlock()

	m.states[service] = authState{
		value:     state,
		expiresAt: time.Now().UTC().Add(10 * time.Minute),
	}

	return state, nil
}

// Checks if the state parameter is valid for the specified service and removes it
func (m *AuthStateStore) Validate(service, state string) error {
	if state == "" {
		return ErrInvalidState
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	entry, exists := m.states[service]
	if !exists || entry.value != state {
		return ErrInvalidState
	}

	// Remove the state parameter (one-time use)
	delete(m.states, service)

	// Check if expired
	if time.Now().UTC().After(entry.expiresAt) {
		return ErrInvalidState
	}

	return nil
}
