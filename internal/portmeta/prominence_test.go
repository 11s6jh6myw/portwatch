package portmeta

import (
	"fmt"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePromPort(port int, meta map[string]string) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Meta: meta}
}

func TestProminenceLevel_String(t *testing.T) {
	cases := []struct {
		level ProminenceLevel
		want  string
	}{
		{ProminenceNone, "none"},
		{ProminenceLow, "low"},
		{ProminenceMedium, "medium"},
		{ProminenceHigh, "high"},
		{ProminenceCritical, "critical"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("ProminenceLevel(%d).String() = %q; want %q", tc.level, got, tc.want)
		}
	}
}

func TestProminenceFor_NoMeta(t *testing.T) {
	p := makePromPort(80, nil)
	if got := ProminenceFor(p); got != ProminenceNone {
		t.Errorf("expected none, got %s", got)
	}
}

func TestProminenceFor_RecentlyActivePort(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)
	p := makePromPort(443, map[string]string{
		"last_seen":  now,
		"prevalent":  "true",
		"critical":   "true",
		"high_impact": "true",
	})
	got := ProminenceFor(p)
	if got < ProminenceHigh {
		t.Errorf("expected at least high, got %s", got)
	}
}

func TestProminenceAnnotator_AddsMetadata(t *testing.T) {
	annotate := NewProminenceAnnotator()
	ports := []scanner.PortInfo{
		makePromPort(22, map[string]string{}),
	}
	out := annotate(ports)
	if _, ok := out[0].Meta["prominence"]; !ok {
		t.Error("expected prominence key in meta")
	}
}

func TestFilterByMinProminence_IncludesHighEnough(t *testing.T) {
	ports := []scanner.PortInfo{
		makePromPort(80, map[string]string{"prominence": "high"}),
		makePromPort(9999, map[string]string{"prominence": "low"}),
		makePromPort(443, map[string]string{"prominence": "critical"}),
	}
	out := FilterByMinProminence(ports, ProminenceHigh)
	if len(out) != 2 {
		t.Errorf("expected 2 ports, got %d", len(out))
	}
}

func TestFilterByMinProminence_NoMeta_PassesThrough(t *testing.T) {
	ports := []scanner.PortInfo{
		makePromPort(1234, nil),
	}
	out := FilterByMinProminence(ports, ProminenceCritical)
	if len(out) != 1 {
		t.Errorf("expected port to pass through without meta, got %d", len(out))
	}
}

func TestParseProminence_RoundTrip(t *testing.T) {
	levels := []ProminenceLevel{ProminenceNone, ProminenceLow, ProminenceMedium, ProminenceHigh, ProminenceCritical}
	for _, l := range levels {
		if got := parseProminence(l.String()); got != l {
			t.Errorf("round-trip failed for %s: got %s", l, fmt.Sprint(got))
		}
	}
}
