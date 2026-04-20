package portmeta_test

import (
	"testing"

	"github.com/iamcalledrob/portwatch/internal/portmeta"
	"github.com/iamcalledrob/portwatch/internal/scanner"
)

func makeImpactPort(port int) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Proto: "tcp"}
}

func TestImpactLevel_String(t *testing.T) {
	cases := []struct {
		level portmeta.ImpactLevel
		want  string
	}{
		{portmeta.ImpactNone, "none"},
		{portmeta.ImpactLow, "low"},
		{portmeta.ImpactMedium, "medium"},
		{portmeta.ImpactHigh, "high"},
		{portmeta.ImpactCritical, "critical"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("ImpactLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestImpactFor_CriticalPort(t *testing.T) {
	// Port 22 (SSH) should be at least high impact.
	p := makeImpactPort(22)
	got := portmeta.ImpactFor(p)
	if got < portmeta.ImpactHigh {
		t.Errorf("ImpactFor(22) = %v, want >= high", got)
	}
}

func TestImpactFor_UnknownPort(t *testing.T) {
	p := makeImpactPort(19999)
	got := portmeta.ImpactFor(p)
	if got != portmeta.ImpactNone {
		t.Errorf("ImpactFor(19999) = %v, want none", got)
	}
}

func TestImpactFor_HTTPPort(t *testing.T) {
	p := makeImpactPort(80)
	got := portmeta.ImpactFor(p)
	if got == portmeta.ImpactNone {
		t.Errorf("ImpactFor(80) = none, want at least low")
	}
}

func TestIsHighImpact_True(t *testing.T) {
	p := makeImpactPort(22)
	if !portmeta.IsHighImpact(p) {
		t.Errorf("IsHighImpact(22) = false, want true")
	}
}

func TestIsHighImpact_False(t *testing.T) {
	p := makeImpactPort(19999)
	if portmeta.IsHighImpact(p) {
		t.Errorf("IsHighImpact(19999) = true, want false")
	}
}
