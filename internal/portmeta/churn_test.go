package portmeta

import (
	"testing"
	"time"
)

func TestChurnLevel_String(t *testing.T) {
	cases := []struct {
		level ChurnLevel
		want  string
	}{
		{ChurnNone, "none"},
		{ChurnLow, "low"},
		{ChurnModerate, "moderate"},
		{ChurnHigh, "high"},
		{ChurnLevel(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("ChurnLevel(%d).String() = %q; want %q", tc.level, got, tc.want)
		}
	}
}

func TestClassifyChurn_NoEvents(t *testing.T) {
	if got := ClassifyChurn(nil, time.Hour); got != ChurnNone {
		t.Errorf("expected ChurnNone, got %s", got)
	}
}

func TestClassifyChurn_LowChurn(t *testing.T) {
	events := []ChurnEvent{
		{At: time.Now().Add(-10 * time.Minute), Kind: "opened"},
		{At: time.Now().Add(-5 * time.Minute), Kind: "closed"},
	}
	if got := ClassifyChurn(events, time.Hour); got != ChurnLow {
		t.Errorf("expected ChurnLow, got %s", got)
	}
}

func TestClassifyChurn_ModerateChurn(t *testing.T) {
	events := make([]ChurnEvent, 4)
	for i := range events {
		events[i] = ChurnEvent{At: time.Now().Add(-time.Duration(i+1) * time.Minute), Kind: "opened"}
	}
	if got := ClassifyChurn(events, time.Hour); got != ChurnModerate {
		t.Errorf("expected ChurnModerate, got %s", got)
	}
}

func TestClassifyChurn_HighChurn(t *testing.T) {
	events := make([]ChurnEvent, 8)
	for i := range events {
		events[i] = ChurnEvent{At: time.Now().Add(-time.Duration(i+1) * time.Minute), Kind: "opened"}
	}
	if got := ClassifyChurn(events, time.Hour); got != ChurnHigh {
		t.Errorf("expected ChurnHigh, got %s", got)
	}
}

func TestClassifyChurn_EventsOutsideWindowIgnored(t *testing.T) {
	events := []ChurnEvent{
		{At: time.Now().Add(-2 * time.Hour), Kind: "opened"},
		{At: time.Now().Add(-3 * time.Hour), Kind: "closed"},
	}
	if got := ClassifyChurn(events, time.Hour); got != ChurnNone {
		t.Errorf("expected ChurnNone for stale events, got %s", got)
	}
}
