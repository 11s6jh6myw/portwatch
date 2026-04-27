package portmeta

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makeDepPort(port int) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Proto: "tcp"}
}

func TestDeprecationAnnotator_AddsMetadata(t *testing.T) {
	annotate := NewDeprecationAnnotator()
	ports := annotate([]scanner.PortInfo{makeDepPort(23)}) // Telnet
	if len(ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(ports))
	}
	p := ports[0]
	if p.Meta == nil {
		t.Fatal("expected Meta to be set")
	}
	if _, ok := p.Meta[metaDeprecationLevel]; !ok {
		t.Error("expected deprecation.level key in Meta")
	}
	if _, ok := p.Meta[metaDeprecationScore]; !ok {
		t.Error("expected deprecation.score key in Meta")
	}
}

func TestDeprecationAnnotator_ActivePort_NoneLevel(t *testing.T) {
	annotate := NewDeprecationAnnotator()
	ports := annotate([]scanner.PortInfo{makeDepPort(443)})
	got := ports[0].Meta[metaDeprecationLevel]
	if got != "none" {
		t.Errorf("expected none for HTTPS, got %q", got)
	}
}

func TestFilterByMaxDeprecation_IncludesNone(t *testing.T) {
	annotate := NewDeprecationAnnotator()
	ports := annotate([]scanner.PortInfo{
		makeDepPort(443), // none
		makeDepPort(23),  // high/obsolete
	})
	filtered := FilterByMaxDeprecation(ports, DeprecationNone)
	if len(filtered) != 1 {
		t.Fatalf("expected 1 port, got %d", len(filtered))
	}
	if filtered[0].Port != 443 {
		t.Errorf("expected port 443, got %d", filtered[0].Port)
	}
}

func TestFilterByMaxDeprecation_NoMeta_PassesThrough(t *testing.T) {
	p := makeDepPort(9999) // no meta
	result := FilterByMaxDeprecation([]scanner.PortInfo{p}, DeprecationNone)
	if len(result) != 1 {
		t.Errorf("expected port without meta to pass through, got %d results", len(result))
	}
}

func TestFilterByMaxDeprecation_IncludesUpToMax(t *testing.T) {
	annotate := NewDeprecationAnnotator()
	ports := annotate([]scanner.PortInfo{
		makeDepPort(443), // none
		makeDepPort(21),  // high/obsolete
		makeDepPort(23),  // high/obsolete
	})
	// Allow up to high — should include none and high but not obsolete (if any)
	filtered := FilterByMaxDeprecation(ports, DeprecationHigh)
	for _, p := range filtered {
		level := parseDeprecation(p.Meta[metaDeprecationLevel])
		if level > DeprecationHigh {
			t.Errorf("port %d has level %v, exceeds max high", p.Port, level)
		}
	}
}
