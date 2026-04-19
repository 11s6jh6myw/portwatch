package portmeta

import (
	"testing"

	"github.com/iamcathal/portwatch/internal/scanner"
)

func TestAnomalyLevel_String(t *testing.T) {
	cases := []struct {
		level AnomalyLevel
		want  string
	}{
		{AnomalyNone, "none"},
		{AnomalyLow, "low"},
		{AnomalyMedium, "medium"},
		{AnomalyHigh, "high"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("AnomalyLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestAnomalyFor_HighRiskPort(t *testing.T) {
	// Port 23 (telnet) should score high anomaly.
	level := AnomalyFor(23)
	if level < AnomalyMedium {
		t.Errorf("expected at least medium anomaly for port 23, got %s", level)
	}
}

func TestAnomalyFor_SafePort(t *testing.T) {
	// Port 443 (HTTPS) should have low or no anomaly.
	level := AnomalyFor(443)
	if level >= AnomalyHigh {
		t.Errorf("expected less than high anomaly for port 443, got %s", level)
	}
}

func TestAnomalyFor_UnknownPort(t *testing.T) {
	level := AnomalyFor(61234)
	if level > AnomalyMedium {
		t.Errorf("unexpected high anomaly for unknown port, got %s", level)
	}
}

func TestFilterByMinAnomaly_FiltersCorrectly(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 443},  // low anomaly
		{Port: 23},   // high anomaly
		{Port: 4444}, // likely medium/high
	}
	result := FilterByMinAnomaly(ports, AnomalyHigh)
	for _, p := range result {
		if AnomalyFor(p.Port) < AnomalyHigh {
			t.Errorf("port %d should not pass AnomalyHigh filter", p.Port)
		}
	}
}

func TestAnomalyAnnotator_AddsMetadata(t *testing.T) {
	annotate := NewAnomalyAnnotator()
	ports := []scanner.PortInfo{{Port: 23}}
	out := annotate(ports)
	if out[0].Metadata == nil {
		t.Fatal("expected metadata to be set")
	}
	if _, ok := out[0].Metadata["anomaly"]; !ok {
		t.Error("expected 'anomaly' key in metadata")
	}
}
