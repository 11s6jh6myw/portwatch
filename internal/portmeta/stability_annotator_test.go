package portmeta

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeStabilityPort(seenAgo time.Duration, changes int) scanner.PortInfo {
	p := scanner.PortInfo{Port: 80, Meta: make(map[string]string)}
	if seenAgo > 0 {
		p.Meta[stabilitySeenKey] = time.Now().Add(-seenAgo).Format(time.RFC3339)
	}
	if changes > 0 {
		p.Meta[stabilityChangesKey] = itoa(changes)
	}
	return p
}

func TestStabilityAnnotator_AddsMetadata(t *testing.T) {
	annotate := NewStabilityAnnotator()
	ports := []scanner.PortInfo{makeStabilityPort(50*24*time.Hour, 0)}
	result := annotate(ports)
	if result[0].Meta[stabilityKey] != "locked" {
		t.Errorf("expected locked, got %q", result[0].Meta[stabilityKey])
	}
}

func TestStabilityAnnotator_UnstablePort(t *testing.T) {
	annotate := NewStabilityAnnotator()
	ports := []scanner.PortInfo{makeStabilityPort(time.Hour, 15)}
	result := annotate(ports)
	if result[0].Meta[stabilityKey] != "unstable" {
		t.Errorf("expected unstable, got %q", result[0].Meta[stabilityKey])
	}
}

func TestFilterByMinStability_IncludesStable(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 80, Meta: map[string]string{stabilityKey: "stable"}},
		{Port: 443, Meta: map[string]string{stabilityKey: "locked"}},
		{Port: 9000, Meta: map[string]string{stabilityKey: "unstable"}},
	}
	result := FilterByMinStability(ports, StabilityStable)
	if len(result) != 2 {
		t.Errorf("expected 2 ports, got %d", len(result))
	}
}

func TestFilterByMinStability_NoMeta_PassesThrough(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 80, Meta: map[string]string{}},
	}
	result := FilterByMinStability(ports, StabilityUnknown)
	if len(result) != 1 {
		t.Errorf("expected 1 port, got %d", len(result))
	}
}
