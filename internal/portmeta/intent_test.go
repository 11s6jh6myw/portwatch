package portmeta_test

import (
	"testing"

	"github.com/user/portwatch/internal/portmeta"
	"github.com/user/portwatch/internal/scanner"
)

func TestIntentFor_KnownPorts(t *testing.T) {
	cases := []struct {
		port int
		want portmeta.IntentLevel
	}{
		{22, portmeta.IntentAdministrative},
		{23, portmeta.IntentLegacy},
		{80, portmeta.IntentApplication},
		{8080, portmeta.IntentDevelopment},
		{53, portmeta.IntentInfrastructure},
	}
	for _, tc := range cases {
		got := portmeta.IntentFor(tc.port)
		if got != tc.want {
			t.Errorf("IntentFor(%d) = %v, want %v", tc.port, got, tc.want)
		}
	}
}

func TestIntentFor_UnknownPort(t *testing.T) {
	if got := portmeta.IntentFor(9999); got != portmeta.IntentUnknown {
		t.Errorf("expected IntentUnknown, got %v", got)
	}
}

func TestIntentLevel_String(t *testing.T) {
	if portmeta.IntentAdministrative.String() != "administrative" {
		t.Error("unexpected string for IntentAdministrative")
	}
	if portmeta.IntentUnknown.String() != "unknown" {
		t.Error("unexpected string for IntentUnknown")
	}
}

func TestIntentAnnotator_AddsMetadata(t *testing.T) {
	annotate := portmeta.NewIntentAnnotator()
	ports := []scanner.PortInfo{{Port: 22}, {Port: 80}, {Port: 9999}}
	result := annotate(ports)
	if result[0].Meta["intent"] != "administrative" {
		t.Errorf("expected administrative, got %s", result[0].Meta["intent"])
	}
	if result[1].Meta["intent"] != "application" {
		t.Errorf("expected application, got %s", result[1].Meta["intent"])
	}
	if result[2].Meta["intent"] != "unknown" {
		t.Errorf("expected unknown, got %s", result[2].Meta["intent"])
	}
}

func TestFilterByIntent_ReturnsMatching(t *testing.T) {
	ports := []scanner.PortInfo{{Port: 22}, {Port: 80}, {Port: 8080}}
	got := portmeta.FilterByIntent(ports, portmeta.IntentDevelopment)
	if len(got) != 1 || got[0].Port != 8080 {
		t.Errorf("expected only port 8080, got %v", got)
	}
}
