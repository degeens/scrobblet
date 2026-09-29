package common

import (
	"sync"
	"time"
)

// HealthStatus tracks the health of a source or target in a thread-safe manner.
type HealthStatus struct {
	mu              sync.RWMutex
	healthy         bool
	lastHealthCheck time.Time
}

// NewHealthStatus creates a HealthStatus initialized to healthy with the current UTC time.
func NewHealthStatus() *HealthStatus {
	return &HealthStatus{
		healthy:         true,
		lastHealthCheck: time.Now().UTC(),
	}
}

// Get returns the current health state and the timestamp of the last health check.
func (s *HealthStatus) Get() (healthy bool, lastHealthCheck time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.healthy, s.lastHealthCheck
}

// Set updates the health state and records the current UTC time as the last health check.
func (s *HealthStatus) Set(healthy bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.healthy = healthy
	s.lastHealthCheck = time.Now().UTC()
}
