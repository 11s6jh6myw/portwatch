package portmeta

import (
	"testing"
)

func TestScore_HighRiskPort(t *testing.T) {
	highRisk := []int{23, 512, 513, 514, 1433, 3389, 5900}
	for _, p := range highRisk {
		if got := Score(p); got != RiskHigh {
			t.Errorf("Score(%d) = %v, want high", p, got)
		}
	}
}

func TestScore_MediumRiskPort(t *testing.T) {
	medium := []int{21, 69, 161, 3306, 5432, 6379, 9200}
	for _, p := range medium {
		if got := Score(p); got != RiskMedium {
			t.Errorf("Score(%d) = %v, want medium", p, got)
		}
	}
}

func TestScore_LowRiskPort(t *testing.T) {
	// IsRisky covers ports like 4444, 1234 — use one known risky but not high/medium
	// We test a port that IsRisky returns true for but isn't in high/medium maps.
	// Port 4444 is a known risky port in portmeta.
	if got := Score(4444); got != RiskLow {
		t.Errorf("Score(4444) = %v, want low", got)
	}
}

func TestScore_NoRiskPort(t *testing.T) {
	safe := []int{80, 443, 8080, 12345}
	for _, p := range safe {
		if got := Score(p); got != RiskNone {
			t.Errorf("Score(%d) = %v, want none", p, got)
		}
	}
}

func TestRiskLevel_String(t *testing.T) {
	cases := []struct {
		level RiskLevel
		want  string
	}{
		{RiskNone, "none"},
		{RiskLow, "low"},
		{RiskMedium, "medium"},
		{RiskHigh, "high"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf"RiskLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}
