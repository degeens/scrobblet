package health

import (
	"time"

	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

type HealthCheck struct {
	Healthy bool
	Source  SourceHealthCheck
	Targets []TargetHealthCheck
}

type SourceHealthCheck struct {
	SourceType      sources.SourceType
	Healthy         bool
	LastHealthCheck time.Time
}

type TargetHealthCheck struct {
	TargetType      targets.TargetType
	Healthy         bool
	LastHealthCheck time.Time
}

func CheckHealth(source sources.Source, targets []targets.Target) HealthCheck {
	sourceCheck, sourceHealthy := checkSource(source)
	targetChecks, targetsHealthy := checkTargets(targets)

	healthy := sourceHealthy && targetsHealthy

	return HealthCheck{
		Healthy: healthy,
		Source:  sourceCheck,
		Targets: targetChecks,
	}
}

func checkSource(source sources.Source) (SourceHealthCheck, bool) {
	healthy, lastHealthCheck := source.Healthy()

	healthCheck := SourceHealthCheck{SourceType: source.SourceType(), Healthy: healthy, LastHealthCheck: lastHealthCheck}

	return healthCheck, healthy
}

func checkTargets(targets []targets.Target) ([]TargetHealthCheck, bool) {
	allHealthy := true

	healthChecks := make([]TargetHealthCheck, 0, len(targets))

	for _, target := range targets {
		healthy, lastHealthCheck := target.Healthy()

		if !healthy {
			allHealthy = false
		}

		healthCheck := TargetHealthCheck{TargetType: target.TargetType(), Healthy: healthy, LastHealthCheck: lastHealthCheck}

		healthChecks = append(healthChecks, healthCheck)
	}

	return healthChecks, allHealthy
}
