package portmeta

import (
	"testing"
)

func TestCriticalityLevel_String(t *testing.T) {
	cases := []struct {
		level CriticalityLevel
		want  string
	}{
		{CriticalityNone, "none"},
		{CriticalityLow, "low"},
		{CriticalityMedium, "medium"},
		{CriticalityHigh, "high"},
		{CriticalityCritical, "critical"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}

func TestCriticalityFor_KnownPorts(t *testing.T) {
	cases := []struct {
		port uint16
		want CriticalityLevel
	}{
		{53, CriticalityCritical},
		{443, CriticalityHigh},
		{80, CriticalityMedium},
		{23, CriticalityLow},
	}
	for _, tc := range cases {
		if got := CriticalityFor(tc.port); got != tc.want {
			t.Errorf("CriticalityFor(%d) = %v, want %v", tc.port, got, tc.want)
		}
	}
}

func TestCriticalityFor_UnknownPort(t *testing.T) {
	if got := CriticalityFor(9999); got != CriticalityNone {
		t.Errorf("expected CriticalityNone for unknown port, got %v", got)
	}
}

func TestIsCritical_True(t *testing.T) {
	for _, port := range []uint16{22, 53, 3306, 6443} {
		if !IsCritical(port) {
			t.Errorf("expected port %d to be critical", port)
		}
	}
}

func TestIsCritical_False(t *testing.T) {
	for _, port := range []uint16{23, 80, 9999} {
		if IsCritical(port) {
			t.Errorf("expected port %d to not be critical", port)
		}
	}
}
