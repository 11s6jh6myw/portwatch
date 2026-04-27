package portmeta

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeRecencyPort(lastSeen time.Time) scanner.PortInfo {
	meta := map[string]string{}
	if !lastSeen.IsZero() {
		meta["last_seen"] = lastSeen.Format(time.RFC3339)
	}
	return scanner.PortInfo{Port: 80, Meta: meta}
}

func TestRecencyLevel_String(t *testing.T) {
	cases := []struct {
		level RecencyLevel
		want  string
	}{
		{RecencyNone, "none"},
		{RecencyStale, "stale"},
		{RecencyAged, "aged"},
		{RecencyRecent, "recent"},
		{RecencyLive, "live"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("RecencyLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestRecencyFor_NoMeta(t *testing.T) {
	p := scanner.PortInfo{Port: 443}
	if got := RecencyFor(p); got != RecencyNone {
		t.Errorf("expected RecencyNone, got %v", got)
	}
}

func TestRecencyFor_Live(t *testing.T) {
	p := makeRecencyPort(time.Now().Add(-2 * time.Minute))
	if got := RecencyFor(p); got != RecencyLive {
		t.Errorf("expected RecencyLive, got %v", got)
	}
}

func TestRecencyFor_Recent(t *testing.T) {
	p := makeRecencyPort(time.Now().Add(-6 * time.Hour))
	if got := RecencyFor(p); got != RecencyRecent {
		t.Errorf("expected RecencyRecent, got %v", got)
	}
}

func TestRecencyFor_Aged(t *testing.T) {
	p := makeRecencyPort(time.Now().Add(-72 * time.Hour))
	if got := RecencyFor(p); got != RecencyAged {
		t.Errorf("expected RecencyAged, got %v", got)
	}
}

func TestRecencyFor_Stale(t *testing.T) {
	p := makeRecencyPort(time.Now().Add(-10 * 24 * time.Hour))
	if got := RecencyFor(p); got != RecencyStale {
		t.Errorf("expected RecencyStale, got %v", got)
	}
}

func TestRecencyAnnotator_AddsMetadata(t *testing.T) {
	annotate := NewRecencyAnnotator()
	ports := []scanner.PortInfo{makeRecencyPort(time.Now().Add(-1 * time.Minute))}
	result := annotate(ports)
	if result[0].Meta["recency"] != "live" {
		t.Errorf("expected recency=live, got %q", result[0].Meta["recency"])
	}
	if result[0].Meta["is_live"] != "true" {
		t.Errorf("expected is_live=true, got %q", result[0].Meta["is_live"])
	}
}

func TestFilterByMinRecency_ExcludesStale(t *testing.T) {
	ports := []scanner.PortInfo{
		makeRecencyPort(time.Now().Add(-20 * 24 * time.Hour)), // stale
		makeRecencyPort(time.Now().Add(-1 * time.Minute)),     // live
	}
	result := FilterByMinRecency(ports, RecencyRecent)
	if len(result) != 1 {
		t.Fatalf("expected 1 port, got %d", len(result))
	}
	if RecencyFor(result[0]) != RecencyLive {
		t.Errorf("expected live port to pass filter")
	}
}
