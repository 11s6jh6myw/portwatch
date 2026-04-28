package portmeta

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makeFreqPort(seen, scans int) scanner.PortInfo {
	p := scanner.PortInfo{Port: 80, Meta: map[string]string{}}
	if seen >= 0 {
		p.Meta["seen_count"] = itoa(seen)
	}
	if scans >= 0 {
		p.Meta["scan_count"] = itoa(scans)
	}
	return p
}

func TestFrequencyLevel_String(t *testing.T) {
	cases := []struct {
		level FrequencyLevel
		want  string
	}{
		{FrequencyNone, "none"},
		{FrequencyRare, "rare"},
		{FrequencyOccasional, "occasional"},
		{FrequencyCommon, "common"},
		{FrequencyAlways, "always"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("FrequencyLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestFrequencyFor_NoMeta(t *testing.T) {
	p := scanner.PortInfo{Port: 443}
	if got := FrequencyFor(p); got != FrequencyNone {
		t.Errorf("expected FrequencyNone for nil meta, got %v", got)
	}
}

func TestFrequencyFor_AlwaysPresent(t *testing.T) {
	p := makeFreqPort(100, 100)
	if got := FrequencyFor(p); got != FrequencyAlways {
		t.Errorf("expected FrequencyAlways, got %v", got)
	}
}

func TestFrequencyFor_Common(t *testing.T) {
	p := makeFreqPort(70, 100)
	if got := FrequencyFor(p); got != FrequencyCommon {
		t.Errorf("expected FrequencyCommon, got %v", got)
	}
}

func TestFrequencyFor_Occasional(t *testing.T) {
	p := makeFreqPort(40, 100)
	if got := FrequencyFor(p); got != FrequencyOccasional {
		t.Errorf("expected FrequencyOccasional, got %v", got)
	}
}

func TestFrequencyFor_Rare(t *testing.T) {
	p := makeFreqPort(5, 100)
	if got := FrequencyFor(p); got != FrequencyRare {
		t.Errorf("expected FrequencyRare, got %v", got)
	}
}

func TestFrequencyFor_ZeroScans(t *testing.T) {
	p := makeFreqPort(10, 0)
	if got := FrequencyFor(p); got != FrequencyNone {
		t.Errorf("expected FrequencyNone for zero scans, got %v", got)
	}
}

func TestIsFrequent_True(t *testing.T) {
	p := makeFreqPort(80, 100)
	if !IsFrequent(p) {
		t.Error("expected IsFrequent to return true")
	}
}

func TestIsFrequent_False(t *testing.T) {
	p := makeFreqPort(10, 100)
	if IsFrequent(p) {
		t.Error("expected IsFrequent to return false")
	}
}
