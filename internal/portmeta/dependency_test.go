package portmeta_test

import (
	"testing"

	"github.com/joshbeard/portwatch/internal/portmeta"
	"github.com/joshbeard/portwatch/internal/scanner"
)

func TestDependenciesFor_KnownPort(t *testing.T) {
	deps := portmeta.DependenciesFor(443)
	if len(deps) == 0 {
		t.Fatal("expected dependencies for port 443")
	}
	found := false
	for _, d := range deps {
		if d.Related == 80 || d.Port == 80 {
			found = true
		}
	}
	if !found {
		t.Error("expected port 80 in dependencies of 443")
	}
}

func TestDependenciesFor_UnknownPort(t *testing.T) {
	deps := portmeta.DependenciesFor(9999)
	if len(deps) != 0 {
		t.Errorf("expected no dependencies, got %d", len(deps))
	}
}

func TestHasDependency_True(t *testing.T) {
	if !portmeta.HasDependency(443, 80) {
		t.Error("expected 443 and 80 to be dependent")
	}
	if !portmeta.HasDependency(80, 443) {
		t.Error("expected reverse lookup to work")
	}
}

func TestHasDependency_False(t *testing.T) {
	if portmeta.HasDependency(9999, 8888) {
		t.Error("expected no dependency between unknown ports")
	}
}

func TestDependencyAnnotator_AddsMetadata(t *testing.T) {
	ports := []scanner.PortInfo{{Port: 443}, {Port: 9999}}
	a := portmeta.NewDependencyAnnotator()
	result := a.Annotate(ports)

	if result[0].Meta == nil || result[0].Meta["dep.count"] == "" {
		t.Error("expected dep.count for port 443")
	}
	if result[1].Meta != nil && result[1].Meta["dep.count"] != "" {
		t.Error("expected no dep.count for port 9999")
	}
}

func TestFilterByHasDependency(t *testing.T) {
	ports := []scanner.PortInfo{{Port: 443}, {Port: 9999}, {Port: 6379}}
	out := portmeta.FilterByHasDependency(ports)
	if len(out) != 2 {
		t.Errorf("expected 2 ports with dependencies, got %d", len(out))
	}
}
