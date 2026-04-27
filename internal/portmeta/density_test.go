package portmeta_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portmeta"
	"github.com/user/portwatch/internal/scanner"
)

func makeDensityPort(port int) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Proto: "tcp"}
}

func recentTimes(n int) []time.Time {
	times := make([]time.Time, n)
	for i := range times {
		times[i] = time.Now().Add(-time.Duration(i) * time.Minute)
	}
	return times
}

func TestDensityLevel_String(t *testing.T) {
	cases := []struct {
		level portmeta.DensityLevel
		want  string
	}{
		{portmeta.DensityNone, "none"},
		{portmeta.DensitySparse, "sparse"},
		{portmeta.DensityModerate, "moderate"},
		{portmeta.DensityDense, "dense"},
		{portmeta.DensitySaturated, "saturated"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("DensityLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestDensityFor_NoEvents(t *testing.T) {
	p := makeDensityPort(80)
	got := portmeta.DensityFor(p, nil, time.Hour, 10)
	if got != portmeta.DensityNone {
		t.Errorf("expected DensityNone, got %s", got)
	}
}

func TestDensityFor_ZeroScans(t *testing.T) {
	p := makeDensityPort(80)
	got := portmeta.DensityFor(p, recentTimes(5), time.Hour, 0)
	if got != portmeta.DensityNone {
		t.Errorf("expected DensityNone for zero scans, got %s", got)
	}
}

func TestDensityFor_Saturated(t *testing.T) {
	p := makeDensityPort(443)
	events := recentTimes(10)
	got := portmeta.DensityFor(p, events, time.Hour, 10)
	if got != portmeta.DensitySaturated {
		t.Errorf("expected DensitySaturated, got %s", got)
	}
}

func TestDensityFor_Sparse(t *testing.T) {
	p := makeDensityPort(9999)
	events := recentTimes(2)
	got := portmeta.DensityFor(p, events, time.Hour, 20)
	if got != portmeta.DensitySparse {
		t.Errorf("expected DensitySparse, got %s", got)
	}
}

func TestDensityFor_EventsOutsideWindowIgnored(t *testing.T) {
	p := makeDensityPort(22)
	old := []time.Time{
		time.Now().Add(-3 * time.Hour),
		time.Now().Add(-4 * time.Hour),
	}
	got := portmeta.DensityFor(p, old, time.Hour, 10)
	if got != portmeta.DensityNone {
		t.Errorf("expected DensityNone for stale events, got %s", got)
	}
}
