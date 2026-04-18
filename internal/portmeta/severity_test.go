package portmeta

import (
	"testing"
)

func TestSeverity_String(t *testing.T) {
	cases := []struct {
		s    Severity
		want string
	}{
		{SeverityNone, "none"},
		{SeverityLow, "low"},
		{SeverityMedium, "medium"},
		{SeverityHigh, "high"},
		{SeverityCritical, "critical"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}

func TestSeverityFor_DatabasePort(t *testing.T) {
	// Port 3306 (MySQL) should be critical — high risk + database category.
	s := SeverityFor(3306)
	if s != SeverityCritical && s != SeverityHigh {
		t.Errorf("SeverityFor(3306) = %s, want critical or high", s)
	}
}

func TestSeverityFor_UnknownPort(t *testing.T) {
	s := SeverityFor(39999)
	if s != SeverityNone && s != SeverityMedium {
		t.Errorf("SeverityFor(39999) = %s, want none or medium", s)
	}
}

func TestSeverityFor_HTTPPort(t *testing.T) {
	// Port 80 is well-known but low risk.
	s := SeverityFor(80)
	if s == SeverityCritical {
		t.Errorf("SeverityFor(80) = %s, should not be critical", s)
	}
}

func TestSeverityFor_TelnetPort(t *testing.T) {
	// Port 23 (Telnet) is high risk.
	s := SeverityFor(23)
	if s < SeverityMedium {
		t.Errorf("SeverityFor(23) = %s, want at least medium", s)
	}
}
