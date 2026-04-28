package portmeta

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestSymmetryLevel_String(t *testing.T) {
	cases := []struct {
		level SymmetryLevel
		want  string
	}{
		{SymmetryNone, "none"},
		{SymmetryLow, "low"},
		{SymmetryMedium, "medium"},
		{SymmetryHigh, "high"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("SymmetryLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestSymmetryFor_NoSiblingInPairs(t *testing.T) {
	p := scanner.PortInfo{Port: 9999}
	got := SymmetryFor(p, []scanner.PortInfo{{Port: 80}, {Port: 443}})
	if got != SymmetryNone {
		t.Errorf("expected SymmetryNone for unknown port, got %s", got)
	}
}

func TestSymmetryFor_SiblingAbsent(t *testing.T) {
	p := scanner.PortInfo{Port: 80}
	// sibling 443 is not in the open set
	got := SymmetryFor(p, []scanner.PortInfo{{Port: 22}, {Port: 8080}})
	if got != SymmetryLow {
		t.Errorf("expected SymmetryLow when sibling absent, got %s", got)
	}
}

func TestSymmetryFor_SiblingPresent_DistantPair(t *testing.T) {
	// Port 3306 <-> 33060: distance = 27754, expect SymmetryLow
	p := scanner.PortInfo{Port: 3306}
	open := []scanner.PortInfo{{Port: 3306}, {Port: 33060}}
	got := SymmetryFor(p, open)
	if got != SymmetryLow {
		t.Errorf("expected SymmetryLow for distant sibling pair, got %s", got)
	}
}

func TestSymmetryFor_SiblingPresent_ClosePair(t *testing.T) {
	// Port 20 <-> 21: distance = 1, expect SymmetryHigh
	p := scanner.PortInfo{Port: 20}
	open := []scanner.PortInfo{{Port: 20}, {Port: 21}}
	got := SymmetryFor(p, open)
	if got != SymmetryHigh {
		t.Errorf("expected SymmetryHigh for adjacent sibling pair, got %s", got)
	}
}

func TestSymmetryFor_HTTP_HTTPS_Pair(t *testing.T) {
	// Port 80 <-> 443: distance = 363, expect SymmetryLow when both present
	p := scanner.PortInfo{Port: 80}
	open := []scanner.PortInfo{{Port: 80}, {Port: 443}}
	got := SymmetryFor(p, open)
	if got == SymmetryNone {
		t.Errorf("expected at least SymmetryLow for HTTP/HTTPS pair, got none")
	}
}

func TestSymmetryFor_EmptyOpenSet(t *testing.T) {
	p := scanner.PortInfo{Port: 443}
	got := SymmetryFor(p, []scanner.PortInfo{})
	if got != SymmetryLow {
		t.Errorf("expected SymmetryLow when open set is empty, got %s", got)
	}
}
