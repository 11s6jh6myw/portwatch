package portmeta

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeIsoPort(port uint16) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Proto: "tcp", Meta: make(map[string]string)}
}

func TestIsolationAnnotator_AddsMetadata(t *testing.T) {
	peers := []scanner.PortInfo{makeIsoPort(9000), makeIsoPort(9001)}
	a := NewIsolationAnnotator(peers, nil)

	p := makeIsoPort(9050)
	result := a.Annotate([]scanner.PortInfo{p})

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if _, ok := result[0].Meta["isolation"]; !ok {
		t.Error("expected 'isolation' key in Meta")
	}
}

func TestIsolationAnnotator_HighIsolation_NoPeers(t *testing.T) {
	a := NewIsolationAnnotator(nil, nil)
	p := makeIsoPort(1234)
	result := a.Annotate([]scanner.PortInfo{p})

	if got := result[0].Meta["isolation"]; got != "high" {
		t.Errorf("expected 'high', got %q", got)
	}
}

func TestIsolationAnnotator_NilMetaInitialised(t *testing.T) {
	a := NewIsolationAnnotator(nil, map[uint16]time.Time{443: time.Now()})
	p := scanner.PortInfo{Port: 443, Proto: "tcp"} // no Meta map
	result := a.Annotate([]scanner.PortInfo{p})

	if result[0].Meta == nil {
		t.Error("expected Meta to be initialised")
	}
}

func TestFilterByMinIsolation_IncludesHighEnough(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 80, Proto: "tcp", Meta: map[string]string{"isolation": "high"}},
		{Port: 443, Proto: "tcp", Meta: map[string]string{"isolation": "low"}},
		{Port: 8080, Proto: "tcp", Meta: map[string]string{"isolation": "medium"}},
	}
	result := FilterByMinIsolation(ports, IsolationMedium)
	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
}

func TestFilterByMinIsolation_NoMeta_PassesThrough(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 22, Proto: "tcp"},
	}
	result := FilterByMinIsolation(ports, IsolationHigh)
	if len(result) != 1 {
		t.Errorf("expected port without meta to pass through, got %d", len(result))
	}
}
