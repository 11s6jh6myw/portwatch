package portmeta

import (
	"testing"
	"time"
)

func TestStabilityLevel_String(t *testing.T) {
	cases := []struct {
		level StabilityLevel
		want  string
	}{
		{StabilityUnknown, "unknown"},
		{StabilityUnstable, "unstable"},
		{StabilityVariable, "variable"},
		{StabilityStable, "stable"},
		{StabilityLocked, "locked"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}

func TestClassifyStability_ZeroTime(t *testing.T) {
	if got := ClassifyStability(time.Time{}, 0); got != StabilityUnknown {
		t.Errorf("expected unknown, got %v", got)
	}
}

func TestClassifyStability_Unstable(t *testing.T) {
	if got := ClassifyStability(time.Now().Add(-time.Hour), 12); got != StabilityUnstable {
		t.Errorf("expected unstable, got %v", got)
	}
}

func TestClassifyStability_Variable(t *testing.T) {
	if got := ClassifyStability(time.Now().Add(-time.Hour), 5); got != StabilityVariable {
		t.Errorf("expected variable, got %v", got)
	}
}

func TestClassifyStability_Locked(t *testing.T) {
	old := time.Now().Add(-40 * 24 * time.Hour)
	if got := ClassifyStability(old, 0); got != StabilityLocked {
		t.Errorf("expected locked, got %v", got)
	}
}

func TestClassifyStability_Stable(t *testing.T) {
	if got := ClassifyStability(time.Now().Add(-time.Hour), 2); got != StabilityStable {
		t.Errorf("expected stable, got %v", got)
	}
}
