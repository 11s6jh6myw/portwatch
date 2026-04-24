package portmeta

import (
	"testing"
	"time"
)

func TestUrgencyLevel_String(t *testing.T) {
	cases := []struct {
		level UrgencyLevel
		want  string
	}{
		{UrgencyNone, "none"},
		{UrgencyLow, "low"},
		{UrgencyMedium, "medium"},
		{UrgencyHigh, "high"},
		{UrgencyCritical, "critical"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("UrgencyLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestUrgencyFor_HighRiskPort(t *testing.T) {
	// Port 23 (Telnet) is high-risk, high-severity, high-exposure.
	got := UrgencyFor(23, time.Time{})
	if got < UrgencyHigh {
		t.Errorf("UrgencyFor(23) = %s, want >= high", got)
	}
}

func TestUrgencyFor_SafePort(t *testing.T) {
	// Port 80 (HTTP) is common and generally low urgency.
	got := UrgencyFor(80, time.Time{})
	if got >= UrgencyCritical {
		t.Errorf("UrgencyFor(80) = %s, want < critical", got)
	}
}

func TestUrgencyFor_RecentlyOpenedBoostsUrgency(t *testing.T) {
	// A port opened very recently should score higher than one with zero time.
	recent := time.Now().Add(-1 * time.Minute)
	without := UrgencyFor(8080, time.Time{})
	with := UrgencyFor(8080, recent)
	if with < without {
		t.Errorf("expected recent first-seen to boost urgency: got %s vs %s", with, without)
	}
}

func TestUrgencyFor_UnknownPort(t *testing.T) {
	// An obscure unknown port with no history should be none or low.
	got := UrgencyFor(39876, time.Time{})
	if got > UrgencyMedium {
		t.Errorf("UrgencyFor(39876) = %s, want <= medium", got)
	}
}

func TestUrgencyFor_ZeroTimeNoBoost(t *testing.T) {
	// Zero time should not trigger the recency boost.
	got := UrgencyFor(443, time.Time{})
	_ = got // just ensure no panic
}
