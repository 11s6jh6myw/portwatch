package portmeta

import (
	"testing"
	"time"
)

func TestVolatilityLevel_String(t *testing.T) {
	cases := []struct {
		level VolatilityLevel
		want  string
	}{
		{VolatilityNone, "none"},
		{VolatilityLow, "low"},
		{VolatilityModerate, "moderate"},
		{VolatilityHigh, "high"},
		{VolatilityLevel(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}

func TestVolatilityFor_NoEvents(t *testing.T) {
	level := VolatilityFor(nil, time.Hour)
	if level != VolatilityNone {
		t.Errorf("expected None for empty events, got %s", level)
	}
}

func TestVolatilityFor_ZeroWindow(t *testing.T) {
	events := []time.Time{time.Now()}
	level := VolatilityFor(events, 0)
	if level != VolatilityNone {
		t.Errorf("expected None for zero window, got %s", level)
	}
}

func TestVolatilityFor_LowVolatility(t *testing.T) {
	now := time.Now()
	events := []time.Time{now.Add(-10 * time.Minute), now.Add(-5 * time.Minute)}
	level := VolatilityFor(events, time.Hour)
	if level != VolatilityLow {
		t.Errorf("expected Low, got %s", level)
	}
}

func TestVolatilityFor_ModerateVolatility(t *testing.T) {
	now := time.Now()
	events := make([]time.Time, 5)
	for i := range events {
		events[i] = now.Add(-time.Duration(i+1) * 5 * time.Minute)
	}
	level := VolatilityFor(events, time.Hour)
	if level != VolatilityModerate {
		t.Errorf("expected Moderate, got %s", level)
	}
}

func TestVolatilityFor_HighVolatility(t *testing.T) {
	now := time.Now()
	events := make([]time.Time, 10)
	for i := range events {
		events[i] = now.Add(-time.Duration(i+1) * 2 * time.Minute)
	}
	level := VolatilityFor(events, time.Hour)
	if level != VolatilityHigh {
		t.Errorf("expected High, got %s", level)
	}
}

func TestVolatilityFor_OldEventsIgnored(t *testing.T) {
	now := time.Now()
	// All events are outside the window
	events := []time.Time{
		now.Add(-3 * time.Hour),
		now.Add(-2 * time.Hour),
		now.Add(-90 * time.Minute),
	}
	level := VolatilityFor(events, time.Hour)
	if level != VolatilityNone {
		t.Errorf("expected None when all events outside window, got %s", level)
	}
}
