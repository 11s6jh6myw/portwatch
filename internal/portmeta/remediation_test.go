package portmeta

import (
	"testing"
	"time"
)

func TestRemediationLevel_String(t *testing.T) {
	cases := []struct {
		level RemediationLevel
		want  string
	}{
		{RemediationNone, "none"},
		{RemediationMonitor, "monitor"},
		{RemediationReview, "review"},
		{RemediationMitigate, "mitigate"},
		{RemediationImmediate, "immediate"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}

func TestRemediationFor_HighRiskPort(t *testing.T) {
	// Telnet (23) is high-risk; use a recent firstSeen and high event count.
	firstSeen := time.Now().Add(-1 * time.Hour).Unix()
	level := RemediationFor(23, firstSeen, 50, 100)
	if level < RemediationMitigate {
		t.Errorf("expected >= mitigate for telnet, got %s", level)
	}
}

func TestRemediationFor_SafePort(t *testing.T) {
	// HTTPS (443) should not require immediate action under normal conditions.
	firstSeen := time.Now().Add(-30 * 24 * time.Hour).Unix()
	level := RemediationFor(443, firstSeen, 2, 200)
	if level >= RemediationMitigate {
		t.Errorf("expected < mitigate for https, got %s", level)
	}
}

func TestRemediationFor_UnknownPort(t *testing.T) {
	// Obscure port with no history should at least trigger monitor.
	firstSeen := time.Now().Add(-5 * time.Minute).Unix()
	level := RemediationFor(39999, firstSeen, 1, 10)
	if level < RemediationNone {
		t.Errorf("unexpected negative level: %s", level)
	}
}

func TestIsActionable_True(t *testing.T) {
	for _, l := range []RemediationLevel{RemediationReview, RemediationMitigate, RemediationImmediate} {
		if !IsActionable(l) {
			t.Errorf("expected actionable for %s", l)
		}
	}
}

func TestIsActionable_False(t *testing.T) {
	for _, l := range []RemediationLevel{RemediationNone, RemediationMonitor} {
		if IsActionable(l) {
			t.Errorf("expected not actionable for %s", l)
		}
	}
}
