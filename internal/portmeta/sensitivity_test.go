package portmeta

import (
	"testing"
)

func TestSensitivityLevel_String(t *testing.T) {
	cases := []struct {
		level SensitivityLevel
		want  string
	}{
		{SensitivityNone, "none"},
		{SensitivityLow, "low"},
		{SensitivityMedium, "medium"},
		{SensitivityHigh, "high"},
		{SensitivityCritical, "critical"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}

func TestSensitivityFor_CriticalPort(t *testing.T) {
	if got := SensitivityFor(23); got != SensitivityCritical {
		t.Errorf("port 23: got %v, want critical", got)
	}
}

func TestSensitivityFor_HighPort(t *testing.T) {
	if got := SensitivityFor(22); got != SensitivityHigh {
		t.Errorf("port 22: got %v, want high", got)
	}
}

func TestSensitivityFor_LowPort(t *testing.T) {
	if got := SensitivityFor(80); got != SensitivityLow {
		t.Errorf("port 80: got %v, want low", got)
	}
}

func TestSensitivityFor_UnknownPort(t *testing.T) {
	if got := SensitivityFor(9999); got != SensitivityNone {
		t.Errorf("port 9999: got %v, want none", got)
	}
}

func TestIsSensitive_True(t *testing.T) {
	if !IsSensitive(3306) {
		t.Error("port 3306 should be sensitive")
	}
}

func TestIsSensitive_False(t *testing.T) {
	if IsSensitive(9999) {
		t.Error("port 9999 should not be sensitive")
	}
}

func TestIsSensitive_LowIsNotSensitive(t *testing.T) {
	if IsSensitive(80) {
		t.Error("port 80 (low) should not meet medium threshold")
	}
}
