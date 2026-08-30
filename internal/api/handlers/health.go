package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/degeens/scrobblet/internal/health"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

const (
	StatusHealthy   = "healthy"
	StatusUnhealthy = "unhealthy"
)

type HealthResponse struct {
	Status  string         `json:"status"`
	Source  ClientHealth   `json:"source"`
	Targets []ClientHealth `json:"targets"`
}

type ClientHealth struct {
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func Health(source sources.Source, targets []targets.Target) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		healthCheck := health.CheckHealth(source, targets)

		response := toHealthResponse(healthCheck)

		w.Header().Set("Content-Type", "application/json")

		if healthCheck.Healthy {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func toHealthResponse(h health.HealthCheck) HealthResponse {
	targets := make([]ClientHealth, len(h.Targets))
	for i, t := range h.Targets {
		targets[i] = ClientHealth{
			Type:      string(t.TargetType),
			Status:    toHealthStatus(t.Healthy),
			Timestamp: t.LastHealthCheck,
		}
	}

	return HealthResponse{
		Status: toHealthStatus(h.Healthy),
		Source: ClientHealth{
			Type:      string(h.Source.SourceType),
			Status:    toHealthStatus(h.Source.Healthy),
			Timestamp: h.Source.LastHealthCheck,
		},
		Targets: targets,
	}
}

func toHealthStatus(healthy bool) string {
	if healthy {
		return StatusHealthy
	}
	return StatusUnhealthy
}
