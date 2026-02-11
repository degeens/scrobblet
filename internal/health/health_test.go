package health

import (
	"testing"
	"time"

	"github.com/degeens/scrobblet/internal/common"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

type mockSource struct {
	healthy         bool
	lastHealthCheck time.Time
	sourceType      sources.SourceType
}

func (m *mockSource) Healthy() (bool, time.Time) {
	return m.healthy, m.lastHealthCheck
}

func (m *mockSource) SourceType() sources.SourceType {
	return m.sourceType
}

func (m *mockSource) GetPlaybackState() (*sources.PlaybackState, error) {
	return nil, nil
}

type mockTarget struct {
	healthy         bool
	lastHealthCheck time.Time
	targetType      targets.TargetType
}

func (m *mockTarget) Healthy() (bool, time.Time) {
	return m.healthy, m.lastHealthCheck
}

func (m *mockTarget) TargetType() targets.TargetType {
	return m.targetType
}

func (m *mockTarget) SubmitPlayingTrack(track *common.Track) error {
	return nil
}

func (m *mockTarget) SubmitPlayedTrack(trackedTrack *common.TrackedTrack) error {
	return nil
}

func TestCheckHealth_AllHealthy(t *testing.T) {
	now := time.Now().UTC()

	source := &mockSource{
		healthy:         true,
		lastHealthCheck: now,
		sourceType:      sources.SourceSpotify,
	}

	target := &mockTarget{
		healthy:         true,
		lastHealthCheck: now,
		targetType:      targets.TargetLastFm,
	}

	result := CheckHealth(source, []targets.Target{target})

	if !result.Healthy {
		t.Error("Expected overall health to be true when source and targets are all healthy")
	}

	assertSourceHealthCheck(t, result.Source, true, sources.SourceSpotify, now)

	if len(result.Targets) != 1 {
		t.Fatalf("Expected 1 target, got %d", len(result.Targets))
	}

	assertTargetHealthCheck(t, result.Targets[0], true, targets.TargetLastFm, now)
}

func TestCheckHealth_UnhealthySource(t *testing.T) {
	now := time.Now().UTC()

	source := &mockSource{
		healthy:         false,
		lastHealthCheck: now,
		sourceType:      sources.SourceSpotify,
	}

	target := &mockTarget{
		healthy:         true,
		lastHealthCheck: now,
		targetType:      targets.TargetLastFm,
	}

	result := CheckHealth(source, []targets.Target{target})

	if result.Healthy {
		t.Error("Expected overall health to be false when source is unhealthy")
	}

	assertSourceHealthCheck(t, result.Source, false, sources.SourceSpotify, now)

	if len(result.Targets) != 1 {
		t.Fatalf("Expected 1 target, got %d", len(result.Targets))
	}

	assertTargetHealthCheck(t, result.Targets[0], true, targets.TargetLastFm, now)
}

func TestCheckHealth_UnhealthyTarget(t *testing.T) {
	now := time.Now().UTC()

	source := &mockSource{
		healthy:         true,
		lastHealthCheck: now,
		sourceType:      sources.SourceSpotify,
	}

	target1 := &mockTarget{
		healthy:         true,
		lastHealthCheck: now,
		targetType:      targets.TargetLastFm,
	}

	target2 := &mockTarget{
		healthy:         false,
		lastHealthCheck: now,
		targetType:      targets.TargetListenBrainz,
	}

	result := CheckHealth(source, []targets.Target{target1, target2})

	if result.Healthy {
		t.Error("Expected overall health to be false when a target is unhealthy")
	}

	assertSourceHealthCheck(t, result.Source, true, sources.SourceSpotify, now)

	if len(result.Targets) != 2 {
		t.Fatalf("Expected 2 targets, got %d", len(result.Targets))
	}

	assertTargetHealthCheck(t, result.Targets[0], true, targets.TargetLastFm, now)
	assertTargetHealthCheck(t, result.Targets[1], false, targets.TargetListenBrainz, now)
}

func TestCheckHealth_NoTargets(t *testing.T) {
	now := time.Now().UTC()

	source := &mockSource{
		healthy:         true,
		lastHealthCheck: now,
		sourceType:      sources.SourceSpotify,
	}

	result := CheckHealth(source, []targets.Target{})

	if !result.Healthy {
		t.Error("Expected overall health to be true when source is healthy and no targets")
	}

	assertSourceHealthCheck(t, result.Source, true, sources.SourceSpotify, now)

	if len(result.Targets) != 0 {
		t.Errorf("Expected 0 targets, got %d", len(result.Targets))
	}
}

func assertSourceHealthCheck(t *testing.T, actualSourceHealthCheck SourceHealthCheck, expectedHealthy bool, expectedSourceType sources.SourceType, expectedLastHealthCheck time.Time) {
	t.Helper()

	if actualSourceHealthCheck.Healthy != expectedHealthy {
		t.Errorf("Expected source health to be %v, got %v", expectedHealthy, actualSourceHealthCheck.Healthy)
	}

	if actualSourceHealthCheck.SourceType != expectedSourceType {
		t.Errorf("Expected source type to be %s, got %s", expectedSourceType, actualSourceHealthCheck.SourceType)
	}

	if actualSourceHealthCheck.LastHealthCheck != expectedLastHealthCheck {
		t.Error("Expected source last health check time to match")
	}
}

func assertTargetHealthCheck(t *testing.T, actualTargetHealthCheck TargetHealthCheck, expectedHealthy bool, expectedTargetType targets.TargetType, expectedLastHealthCheck time.Time) {
	t.Helper()

	if actualTargetHealthCheck.Healthy != expectedHealthy {
		t.Errorf("Expected target health to be %v, got %v", expectedHealthy, actualTargetHealthCheck.Healthy)
	}

	if actualTargetHealthCheck.TargetType != expectedTargetType {
		t.Errorf("Expected target type to be %s, got %s", expectedTargetType, actualTargetHealthCheck.TargetType)
	}

	if actualTargetHealthCheck.LastHealthCheck != expectedLastHealthCheck {
		t.Error("Expected target last health check time to match")
	}
}
