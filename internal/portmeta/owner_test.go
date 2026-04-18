package portmeta_test

import (
	"testing"

	"github.com/user/portwatch/internal/portmeta"
	"github.com/user/portwatch/internal/scanner"
)

func TestLookupOwner_KnownPort(t *testing.T) {
	o, ok := portmeta.LookupOwner(443)
	if !ok {
		t.Fatal("expected owner for port 443")
	}
	if o.Org == "" {
		t.Error("expected non-empty Org")
	}
	if o.Contact == "" {
		t.Error("expected non-empty Contact")
	}
}

func TestLookupOwner_UnknownPort(t *testing.T) {
	_, ok := portmeta.LookupOwner(9999)
	if ok {
		t.Error("expected no owner for port 9999")
	}
}

func TestKnownOwner_True(t *testing.T) {
	if !portmeta.KnownOwner(80) {
		t.Error("expected port 80 to have a known owner")
	}
}

func TestKnownOwner_False(t *testing.T) {
	if portmeta.KnownOwner(12345) {
		t.Error("expected port 12345 to have no known owner")
	}
}

func TestOwnerAnnotator_AddsMetadata(t *testing.T) {
	ports := []scanner.PortInfo{{Port: 3306}, {Port: 9999}}
	a := portmeta.NewOwnerAnnotator()
	result := a.Annotate(ports)

	if result[0].Meta["owner_org"] != "Oracle" {
		t.Errorf("expected Oracle, got %q", result[0].Meta["owner_org"])
	}
	if result[1].Meta["owner_org"] != "" {
		t.Errorf("expected empty org for unknown port, got %q", result[1].Meta["owner_org"])
	}
}

func TestFilterByKnownOwner(t *testing.T) {
	ports := []scanner.PortInfo{{Port: 22}, {Port: 8888}, {Port: 5432}}
	result := portmeta.FilterByKnownOwner(ports)
	if len(result) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(result))
	}
	if result[0].Port != 22 || result[1].Port != 5432 {
		t.Errorf("unexpected ports: %v", result)
	}
}
