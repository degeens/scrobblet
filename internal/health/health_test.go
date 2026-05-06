package health

import (
	"testing"

	"github.com/degeens/scrobblet/internal/common"
	"github.com/degeens/scrobblet/internal/sources"
	"github.com/degeens/scrobblet/internal/targets"
)

type mockSource struct {
	healthy    bool
	sourceType sources.SourceType
}

func (m *mockSource) Healthy() bool {
	return m.healthy
}

func (m *mockSource) SourceType() sources.SourceType {
	return m.sourceType
}

func (m *mockSource) GetPlaybackState() (*sources.PlaybackState, error) {
	return nil, nil
}

type mockTarget struct {
	healthy    bool
	targetType targets.TargetType
}

func (m *mockTarget) Healthy() bool {
	return m.healthy
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
	source := &mockSource{
		healthy:    true,
		sourceType: sources.SourceSpotify,
	}

	target := &mockTarget{
		healthy:    true,
		targetType: targets.TargetLastFm,
	}

	result := CheckHealth(source, []targets.Target{target})

	if !result.Healthy {
		t.Error("Expected overall health to be true when source and targets are all healthy")
	}

	assertSourceHealthCheck(t, result.Source, true, sources.SourceSpotify)

	if len(result.Targets) != 1 {
		t.Fatalf("Expected 1 target, got %d", len(result.Targets))
	}

	assertTargetHealthCheck(t, result.Targets[0], true, targets.TargetLastFm)
}

func TestCheckHealth_UnhealthySource(t *testing.T) {
	source := &mockSource{
		healthy:    false,
		sourceType: sources.SourceSpotify,
	}

	target := &mockTarget{
		healthy:    true,
		targetType: targets.TargetLastFm,
	}

	result := CheckHealth(source, []targets.Target{target})

	if result.Healthy {
		t.Error("Expected overall health to be false when source is unhealthy")
	}

	assertSourceHealthCheck(t, result.Source, false, sources.SourceSpotify)

	if len(result.Targets) != 1 {
		t.Fatalf("Expected 1 target, got %d", len(result.Targets))
	}

	assertTargetHealthCheck(t, result.Targets[0], true, targets.TargetLastFm)
}

func TestCheckHealth_UnhealthyTarget(t *testing.T) {
	source := &mockSource{
		healthy:    true,
		sourceType: sources.SourceSpotify,
	}

	target1 := &mockTarget{
		healthy:    true,
		targetType: targets.TargetLastFm,
	}

	target2 := &mockTarget{
		healthy:    false,
		targetType: targets.TargetListenBrainz,
	}

	result := CheckHealth(source, []targets.Target{target1, target2})

	if result.Healthy {
		t.Error("Expected overall health to be false when a target is unhealthy")
	}

	assertSourceHealthCheck(t, result.Source, true, sources.SourceSpotify)

	if len(result.Targets) != 2 {
		t.Fatalf("Expected 2 targets, got %d", len(result.Targets))
	}

	assertTargetHealthCheck(t, result.Targets[0], true, targets.TargetLastFm)
	assertTargetHealthCheck(t, result.Targets[1], false, targets.TargetListenBrainz)
}

func TestCheckHealth_NoTargets(t *testing.T) {
	source := &mockSource{
		healthy:    true,
		sourceType: sources.SourceSpotify,
	}

	result := CheckHealth(source, []targets.Target{})

	if !result.Healthy {
		t.Error("Expected overall health to be true when source is healthy and no targets")
	}

	assertSourceHealthCheck(t, result.Source, true, sources.SourceSpotify)

	if len(result.Targets) != 0 {
		t.Errorf("Expected 0 targets, got %d", len(result.Targets))
	}
}

func assertSourceHealthCheck(t *testing.T, actualSourceHealthCheck SourceHealthCheck, expectedHealthy bool, expectedSourceType sources.SourceType) {
	t.Helper()

	if actualSourceHealthCheck.Healthy != expectedHealthy {
		t.Errorf("Expected source health to be %v, got %v", expectedHealthy, actualSourceHealthCheck.Healthy)
	}

	if actualSourceHealthCheck.SourceType != expectedSourceType {
		t.Errorf("Expected source type to be %s, got %s", expectedSourceType, actualSourceHealthCheck.SourceType)
	}
}

func assertTargetHealthCheck(t *testing.T, actualTargetHealthCheck TargetHealthCheck, expectedHealthy bool, expectedTargetType targets.TargetType) {
	t.Helper()

	if actualTargetHealthCheck.Healthy != expectedHealthy {
		t.Errorf("Expected target health to be %v, got %v", expectedHealthy, actualTargetHealthCheck.Healthy)
	}

	if actualTargetHealthCheck.TargetType != expectedTargetType {
		t.Errorf("Expected target type to be %s, got %s", expectedTargetType, actualTargetHealthCheck.TargetType)
	}
}
