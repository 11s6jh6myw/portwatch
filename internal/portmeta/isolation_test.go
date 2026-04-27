package portmeta

import (
	"testing"

	"github.com/iamcalledrob/portwatch/internal/scanner"
)

func makeIsoPort(port int) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Meta: make(map[string]string)}
}

func TestIsolationLevel_String(t *testing.T) {
	cases := []struct {
		level IsolationLevel
		want  string
	}{
		{IsolationNone, "none"},
		{IsolationLow, "low"},
		{IsolationMedium, "medium"},
		{IsolationHigh, "high"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("IsolationLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestIsolationFor_NoPeers(t *testing.T) {
	p := makeIsoPort(9999)
	got := IsolationFor(p, nil)
	if got != IsolationHigh {
		t.Errorf("expected IsolationHigh with no peers, got %s", got)
	}
}

func TestIsolationFor_ManyKnownNeighbours(t *testing.T) {
	p := makeIsoPort(9999)
	peers := []scanner.PortInfo{
		makeIsoPort(80),
		makeIsoPort(443),
		makeIsoPort(22),
		makeIsoPort(25),
		makeIsoPort(53),
		makeIsoPort(3306),
	}
	got := IsolationFor(p, peers)
	if got != IsolationNone {
		t.Errorf("expected IsolationNone with many known neighbours, got %s", got)
	}
}

func TestIsolationFor_FewKnownNeighbours(t *testing.T) {
	p := makeIsoPort(9999)
	peers := []scanner.PortInfo{makeIsoPort(80), makeIsoPort(443)}
	got := IsolationFor(p, peers)
	if got != IsolationMedium {
		t.Errorf("expected IsolationMedium with 2 known neighbours, got %s", got)
	}
}

func TestIsIsolated_True(t *testing.T) {
	p := makeIsoPort(9999)
	if !IsIsolated(p, nil) {
		t.Error("expected IsIsolated=true with no peers")
	}
}

func TestIsIsolated_False(t *testing.T) {
	p := makeIsoPort(9999)
	peers := []scanner.PortInfo{
		makeIsoPort(80), makeIsoPort(443), makeIsoPort(22),
		makeIsoPort(25), makeIsoPort(53), makeIsoPort(3306),
	}
	if IsIsolated(p, peers) {
		t.Error("expected IsIsolated=false with many known neighbours")
	}
}

func TestIsolationAnnotator_AddsMetadata(t *testing.T) {
	peers := []scanner.PortInfo{makeIsoPort(80)}
	annotate := NewIsolationAnnotator(peers)
	ports := []scanner.PortInfo{makeIsoPort(9999)}
	out := annotate(ports)
	if out[0].Meta[isolationKey] == "" {
		t.Error("expected isolation metadata to be set")
	}
}

func TestFilterByMinIsolation_IncludesMedium(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 1, Meta: map[string]string{isolationKey: "high"}},
		{Port: 2, Meta: map[string]string{isolationKey: "medium"}},
		{Port: 3, Meta: map[string]string{isolationKey: "low"}},
		{Port: 4, Meta: map[string]string{isolationKey: "none"}},
	}
	out := FilterByMinIsolation(ports, IsolationMedium)
	if len(out) != 2 {
		t.Errorf("expected 2 ports, got %d", len(out))
	}
}
