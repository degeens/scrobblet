package health

import (
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

type HealthCheck struct {
	Healthy bool
	Source  SourceHealthCheck
	Targets []TargetHealthCheck
}

type SourceHealthCheck struct {
	SourceType sources.SourceType
	Healthy    bool
}

type TargetHealthCheck struct {
	TargetType targets.TargetType
	Healthy    bool
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
	healthy := source.Healthy()

	healthCheck := SourceHealthCheck{SourceType: source.SourceType(), Healthy: healthy}

	return healthCheck, healthy
}

func checkTargets(targets []targets.Target) ([]TargetHealthCheck, bool) {
	allHealthy := true

	healthChecks := make([]TargetHealthCheck, 0, len(targets))

	for _, target := range targets {
		healthy := target.Healthy()

		if !healthy {
			allHealthy = false
		}

		healthCheck := TargetHealthCheck{TargetType: target.TargetType(), Healthy: healthy}

		healthChecks = append(healthChecks, healthCheck)
	}

	return healthChecks, allHealthy
}
