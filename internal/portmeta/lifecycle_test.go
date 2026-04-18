package portmeta

import (
	"testing"
	"time"
)

func TestClassify_NewPort(t *testing.T) {
	now := time.Now()
	e := LifecycleEntry{Port: 8080, OpenCount: 1, CloseCount: 0, FirstSeen: now.Add(-10 * time.Second), LastSeen: now}
	if got := Classify(e, now, time.Minute); got != StageNew {
		t.Fatalf("expected new, got %s", got)
	}
}

func TestClassify_StablePort(t *testing.T) {
	now := time.Now()
	e := LifecycleEntry{Port: 80, OpenCount: 1, CloseCount: 0, FirstSeen: now.Add(-2 * time.Minute), LastSeen: now}
	if got := Classify(e, now, time.Minute); got != StageStable {
		t.Fatalf("expected stable, got %s", got)
	}
}

func TestClassify_FlappingPort(t *testing.T) {
	now := time.Now()
	e := LifecycleEntry{Port: 9000, OpenCount: 3, CloseCount: 2, FirstSeen: now.Add(-5 * time.Minute), LastSeen: now}
	if got := Classify(e, now, time.Minute); got != StageFlapping {
		t.Fatalf("expected flapping, got %s", got)
	}
}

func TestClassify_RetiredPort(t *testing.T) {
	now := time.Now()
	e := LifecycleEntry{
		Port: 3306, OpenCount: 0, CloseCount: 1,
		FirstSeen: now.Add(-10 * time.Minute),
		LastSeen:  now.Add(-5 * time.Minute),
	}
	if got := Classify(e, now, time.Minute); got != StageRetired {
		t.Fatalf("expected retired, got %s", got)
	}
}

func TestLifecycleStage_String(t *testing.T) {
	for _, tc := range []struct {
		stage LifecycleStage
		want  string
	}{
		{StageNew, "new"},
		{StageStable, "stable"},
		{StageFlapping, "flapping"},
		{StageRetired, "retired"},
	} {
		if tc.stage.String() != tc.want {
			t.Errorf("got %q want %q", tc.stage.String(), tc.want)
		}
	}
}
