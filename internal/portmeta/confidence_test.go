package portmeta

import (
	"testing"
	"time"
)

func TestConfidenceLevel_String(t *testing.T) {
	cases := []struct {
		level ConfidenceLevel
		want  string
	}{
		{ConfidenceLow, "low"},
		{ConfidenceMedium, "medium"},
		{ConfidenceHigh, "high"},
		{ConfidenceCertain, "certain"},
		{ConfidenceLevel(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}

func TestConfidenceFor_ZeroTime(t *testing.T) {
	got := ConfidenceFor(time.Time{}, 0, ChangeFreqStable)
	if got != ConfidenceLow {
		t.Errorf("expected ConfidenceLow, got %v", got)
	}
}

func TestConfidenceFor_VolatileIsLow(t *testing.T) {
	got := ConfidenceFor(time.Now().Add(-30*24*time.Hour), 100, ChangeFreqVolatile)
	if got != ConfidenceLow {
		t.Errorf("expected ConfidenceLow for volatile port, got %v", got)
	}
}

func TestConfidenceFor_Certain(t *testing.T) {
	firstSeen := time.Now().Add(-10 * 24 * time.Hour)
	got := ConfidenceFor(firstSeen, 60, ChangeFreqStable)
	if got != ConfidenceCertain {
		t.Errorf("expected ConfidenceCertain, got %v", got)
	}
}

func TestConfidenceFor_High(t *testing.T) {
	firstSeen := time.Now().Add(-48 * time.Hour)
	got := ConfidenceFor(firstSeen, 25, ChangeFreqOccasional)
	if got != ConfidenceHigh {
		t.Errorf("expected ConfidenceHigh, got %v", got)
	}
}

func TestConfidenceFor_Medium(t *testing.T) {
	firstSeen := time.Now().Add(-2 * time.Hour)
	got := ConfidenceFor(firstSeen, 8, ChangeFreqStable)
	if got != ConfidenceMedium {
		t.Errorf("expected ConfidenceMedium, got %v", got)
	}
}

func TestConfidenceFor_LowFewObservations(t *testing.T) {
	firstSeen := time.Now().Add(-5 * time.Minute)
	got := ConfidenceFor(firstSeen, 2, ChangeFreqStable)
	if got != ConfidenceLow {
		t.Errorf("expected ConfidenceLow, got %v", got)
	}
}
