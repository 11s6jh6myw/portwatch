package portmeta

import (
	"testing"

	"github.com/joeshaw/portwatch/internal/scanner"
)

func TestVisibilityLevel_String(t *testing.T) {
	cases := []struct {
		level VisibilityLevel
		want  string
	}{
		{VisibilityPublic, "public"},
		{VisibilityRestricted, "restricted"},
		{VisibilityInternal, "internal"},
		{VisibilityUnknown, "unknown"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}

func TestVisibilityFor_PublicPort(t *testing.T) {
	if got := VisibilityFor(80); got != VisibilityPublic {
		t.Errorf("expected Public, got %s", got)
	}
}

func TestVisibilityFor_RestrictedPort(t *testing.T) {
	if got := VisibilityFor(22); got != VisibilityRestricted {
		t.Errorf("expected Restricted, got %s", got)
	}
}

func TestVisibilityFor_InternalPort(t *testing.T) {
	if got := VisibilityFor(60000); got != VisibilityInternal {
		t.Errorf("expected Internal, got %s", got)
	}
}

func TestVisibilityFor_UnknownPort(t *testing.T) {
	if got := VisibilityFor(9999); got != VisibilityUnknown {
		t.Errorf("expected Unknown, got %s", got)
	}
}

func TestVisibilityAnnotator_AddsMetadata(t *testing.T) {
	annotate := NewVisibilityAnnotator()
	ports := []scanner.PortInfo{{Port: 443}}
	out := annotate(ports)
	if out[0].Meta["visibility"] != "public" {
		t.Errorf("expected public, got %s", out[0].Meta["visibility"])
	}
}

func TestFilterByMinVisibility_FiltersCorrectly(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 80},    // public
		{Port: 22},    // restricted
		{Port: 60000}, // internal
		{Port: 9999},  // unknown
	}
	out := FilterByMinVisibility(ports, VisibilityRestricted)
	if len(out) != 2 {
		t.Errorf("expected 2 ports, got %d", len(out))
	}
}
