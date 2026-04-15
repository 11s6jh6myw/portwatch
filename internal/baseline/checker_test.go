package baseline_test

import (
	"testing"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/scanner"
)

func snap(ports ...scanner.PortInfo) *baseline.Snapshot {
	return &baseline.Snapshot{Ports: ports}
}

func port(n int) scanner.PortInfo {
	return scanner.PortInfo{Port: n, Proto: "tcp", State: "open"}
}

func TestCheck_NoViolations(t *testing.T) {
	s := snap(port(22), port(443))
	current := []scanner.PortInfo{port(22), port(443)}
	violations := baseline.Check(s, current)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d: %v", len(violations), violations)
	}
}

func TestCheck_UnexpectedOpenPort(t *testing.T) {
	s := snap(port(22))
	current := []scanner.PortInfo{port(22), port(8080)}
	violations := baseline.Check(s, current)
	if len(violations) != 1 {
alf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Port.Port != 8080 {
		t.Errorf("expected violation on port 8080, got %d", violations[0].Port.Port)
	}
}

func TestCheck_MissingBaselinePort(t *testing.T) {
	s := snap(port(22), port(443))
	current := []scanner.PortInfo{port(22)}
	violations := baseline.Check(s, current)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Port.Port != 443 {
		t.Errorf("expected violation on port 443, got %d", violations[0].Port.Port)
	}
}

func TestCheck_EmptyBaseline(t *testing.T) {
	s := snap()
	current := []scanner.PortInfo{port(80)}
	violations := baseline.Check(s, current)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestViolation_String(t *testing.T) {
	v := baseline.Violation{Port: port(22), Reason: "not in baseline"}
	s := v.String()
	if s == "" {
		t.Error("Violation.String() returned empty string")
	}
}
