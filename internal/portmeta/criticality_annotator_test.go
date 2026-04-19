package portmeta

import (
	"testing"

	"github.com/netwatch/portwatch/internal/scanner"
)

func makeCritPort(port int) scanner.PortInfo {
	return scanner.PortInfo{Port: port}
}

func TestCriticalityAnnotator_AddsMetadata(t *testing.T) {
	annotate := NewCriticalityAnnotator()
	ports := []scanner.PortInfo{makeCritPort(53), makeCritPort(9999)}
	out := annotate(ports)

	if out[0].Meta[metaKeyCriticality] != "critical" {
		t.Errorf("expected critical for port 53, got %q", out[0].Meta[metaKeyCriticality])
	}
	if out[1].Meta[metaKeyCriticality] != "none" {
		t.Errorf("expected none for port 9999, got %q", out[1].Meta[metaKeyCriticality])
	}
}

func TestCriticalityAnnotator_ScorePresent(t *testing.T) {
	annotate := NewCriticalityAnnotator()
	out := annotate([]scanner.PortInfo{makeCritPort(443)})
	if out[0].Meta[metaKeyCriticalityScore] == "" {
		t.Error("expected criticality_score to be set")
	}
}

func TestFilterByMinCriticality_IncludesHigh(t *testing.T) {
	ports := []scanner.PortInfo{makeCritPort(53), makeCritPort(443), makeCritPort(80), makeCritPort(9999)}
	out := FilterByMinCriticality(ports, CriticalityHigh)
	if len(out) != 2 {
		t.Errorf("expected 2 ports, got %d", len(out))
	}
}

func TestFilterByMinCriticality_NoMeta_PassesThrough(t *testing.T) {
	ports := []scanner.PortInfo{makeCritPort(22), makeCritPort(6443)}
	out := FilterByMinCriticality(ports, CriticalityCritical)
	if len(out) != 1 {
		t.Errorf("expected 1 critical port, got %d", len(out))
	}
	if out[0].Port != 6443 {
		t.Errorf("expected port 6443, got %d", out[0].Port)
	}
}
