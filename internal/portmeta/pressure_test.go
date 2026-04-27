package portmeta

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePressurePort(port int) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Meta: make(map[string]string)}
}

func recentPressureTimes(n int) []time.Time {
	times := make([]time.Time, n)
	for i := range times {
		times[i] = time.Now().Add(-time.Minute)
	}
	return times
}

func TestPressureLevel_String(t *testing.T) {
	cases := []struct {
		level PressureLevel
		want  string
	}{
		{PressureNone, "none"},
		{PressureLow, "low"},
		{PressureModerate, "moderate"},
		{PressureHigh, "high"},
		{PressureCritical, "critical"},
	}
	for _, c := range cases {
		if got := c.level.String(); got != c.want {
			t.Errorf("PressureLevel(%d).String() = %q, want %q", c.level, got, c.want)
		}
	}
}

func TestPressureFor_NoEvents(t *testing.T) {
	p := makePressurePort(80)
	if got := PressureFor(p, nil, time.Hour); got != PressureNone {
		t.Errorf("expected None, got %v", got)
	}
}

func TestPressureFor_LowCount(t *testing.T) {
	p := makePressurePort(443)
	if got := PressureFor(p, recentPressureTimes(3), time.Hour); got != PressureLow {
		t.Errorf("expected Low, got %v", got)
	}
}

func TestPressureFor_ModerateCount(t *testing.T) {
	p := makePressurePort(22)
	if got := PressureFor(p, recentPressureTimes(10), time.Hour); got != PressureModerate {
		t.Errorf("expected Moderate, got %v", got)
	}
}

func TestPressureFor_HighCount(t *testing.T) {
	p := makePressurePort(3306)
	if got := PressureFor(p, recentPressureTimes(30), time.Hour); got != PressureHigh {
		t.Errorf("expected High, got %v", got)
	}
}

func TestPressureFor_CriticalCount(t *testing.T) {
	p := makePressurePort(23)
	if got := PressureFor(p, recentPressureTimes(60), time.Hour); got != PressureCritical {
		t.Errorf("expected Critical, got %v", got)
	}
}

func TestPressureAnnotator_AddsMetadata(t *testing.T) {
	events := map[int][]time.Time{8080: recentPressureTimes(7)}
	annotate := NewPressureAnnotator(events, time.Hour)
	ports := annotate([]scanner.PortInfo{makePressurePort(8080)})
	if len(ports) != 1 {
		t.Fatal("expected 1 port")
	}
	if got := ports[0].Meta[metaKeyPressure]; got != "moderate" {
		t.Errorf("pressure = %q, want moderate", got)
	}
}

func TestFilterByMaxPressure_IncludesLow(t *testing.T) {
	events := map[int][]time.Time{9000: recentPressureTimes(2), 9001: recentPressureTimes(25)}
	annotate := NewPressureAnnotator(events, time.Hour)
	ports := annotate([]scanner.PortInfo{makePressurePort(9000), makePressurePort(9001)})
	filtered := FilterByMaxPressure(ports, PressureLow)
	if len(filtered) != 1 || filtered[0].Port != 9000 {
		t.Errorf("expected only port 9000, got %v", filtered)
	}
}

func TestFilterByMaxPressure_NoMeta_PassesThrough(t *testing.T) {
	p := scanner.PortInfo{Port: 1234}
	result := FilterByMaxPressure([]scanner.PortInfo{p}, PressureNone)
	if len(result) != 1 {
		t.Error("expected port without meta to pass through")
	}
}
