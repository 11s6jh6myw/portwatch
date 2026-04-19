package portmeta

import (
	"testing"

	"github.com/iamralch/portwatch/internal/scanner"
)

func TestReputationLevel_String(t *testing.T) {
	cases := []struct {
		level ReputationLevel
		want  string
	}{
		{ReputationUnknown, "unknown"},
		{ReputationPoor, "poor"},
		{ReputationFair, "fair"},
		{ReputationGood, "good"},
		{ReputationTrusted, "trusted"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}

func TestReputationFor_KnownPorts(t *testing.T) {
	if got := ReputationFor(443); got != ReputationTrusted {
		t.Errorf("port 443: got %v, want trusted", got)
	}
	if got := ReputationFor(23); got != ReputationPoor {
		t.Errorf("port 23: got %v, want poor", got)
	}
}

func TestReputationFor_UnknownPort(t *testing.T) {
	if got := ReputationFor(9999); got != ReputationUnknown {
		t.Errorf("expected unknown, got %v", got)
	}
}

func TestIsReputable_True(t *testing.T) {
	if !IsReputable(80) {
		t.Error("port 80 should be reputable")
	}
}

func TestIsReputable_False(t *testing.T) {
	if IsReputable(31337) {
		t.Error("port 31337 should not be reputable")
	}
}

func TestReputationAnnotator_AddsMetadata(t *testing.T) {
	ports := []scanner.PortInfo{{Port: 443}, {Port: 23}}
	a := NewReputationAnnotator()
	out := a.Annotate(ports)
	if out[0].Meta["reputation"] != "trusted" {
		t.Errorf("expected trusted, got %s", out[0].Meta["reputation"])
	}
	if out[1].Meta["reputable"] != "false" {
		t.Errorf("expected reputable=false for port 23")
	}
}

func TestFilterByMinReputation_FiltersCorrectly(t *testing.T) {
	ports := []scanner.PortInfo{{Port: 443}, {Port: 23}, {Port: 9999}}
	out := FilterByMinReputation(ports, ReputationGood)
	if len(out) != 1 || out[0].Port != 443 {
		t.Errorf("expected only port 443, got %v", out)
	}
}
