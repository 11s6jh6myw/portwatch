package portmeta

import (
	"testing"
	"time"
)

func TestVelocityLevel_String(t *testing.T) {
	cases := []struct {
		v    VelocityLevel
		want string
	}{
		{VelocityNone, "none"},
		{VelocitySlow, "slow"},
		{VelocityModerate, "moderate"},
		{VelocityFast, "fast"},
		{VelocityRapid, "rapid"},
	}
	for _, c := range cases {
		if got := c.v.String(); got != c.want {
			t.Errorf("VelocityLevel(%d).String() = %q, want %q", c.v, got, c.want)
		}
	}
}

func TestVelocityFor_NoEvents(t *testing.T) {
	if got := VelocityFor(nil, time.Hour); got != VelocityNone {
		t.Errorf("expected None, got %v", got)
	}
}

func TestVelocityFor_SlowVelocity(t *testing.T) {
	events := []time.Time{time.Now().Add(-5 * time.Minute)}
	if got := VelocityFor(events, time.Hour); got != VelocitySlow {
		t.Errorf("expected Slow, got %v", got)
	}
}

func TestVelocityFor_ModerateVelocity(t *testing.T) {
	now := time.Now()
	events := make([]time.Time, 4)
	for i := range events {
		events[i] = now.Add(-time.Duration(i+1) * time.Minute)
	}
	if got := VelocityFor(events, time.Hour); got != VelocityModerate {
		t.Errorf("expected Moderate, got %v", got)
	}
}

func TestVelocityFor_RapidVelocity(t *testing.T) {
	now := time.Now()
	events := make([]time.Time, 15)
	for i := range events {
		events[i] = now.Add(-time.Duration(i+1) * time.Minute)
	}
	if got := VelocityFor(events, time.Hour); got != VelocityRapid {
		t.Errorf("expected Rapid, got %v", got)
	}
}

func TestVelocityFor_EventsOutsideWindowIgnored(t *testing.T) {
	events := []time.Time{
		time.Now().Add(-2 * time.Hour),
		time.Now().Add(-3 * time.Hour),
	}
	if got := VelocityFor(events, time.Hour); got != VelocityNone {
		t.Errorf("expected None for stale events, got %v", got)
	}
}
