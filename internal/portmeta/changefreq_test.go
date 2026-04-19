package portmeta

import (
	"testing"
	"time"
)

func TestChangeFreq_String(t *testing.T) {
	cases := []struct {
		cf   ChangeFreq
		want string
	}{
		{ChangeFreqStable, "stable"},
		{ChangeFreqOccasional, "occasional"},
		{ChangeFreqFrequent, "frequent"},
		{ChangeFreqVolatile, "volatile"},
		{ChangeFreq(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.cf.String(); got != tc.want {
			t.Errorf("ChangeFreq(%d).String() = %q, want %q", tc.cf, got, tc.want)
		}
	}
}

func TestClassifyChangeFreq_Stable(t *testing.T) {
	got := ClassifyChangeFreq(1, 24*time.Hour)
	if got != ChangeFreqStable {
		t.Errorf("expected stable, got %s", got)
	}
}

func TestClassifyChangeFreq_Occasional(t *testing.T) {
	got := ClassifyChangeFreq(3, 2*time.Hour) // 1.5/hr
	if got != ChangeFreqOccasional {
		t.Errorf("expected occasional, got %s", got)
	}
}

func TestClassifyChangeFreq_Frequent(t *testing.T) {
	got := ClassifyChangeFreq(12, time.Hour) // 12/hr
	if got != ChangeFreqFrequent {
		t.Errorf("expected frequent, got %s", got)
	}
}

func TestClassifyChangeFreq_Volatile(t *testing.T) {
	got := ClassifyChangeFreq(40, time.Hour) // 40/hr
	if got != ChangeFreqVolatile {
		t.Errorf("expected volatile, got %s", got)
	}
}

func TestClassifyChangeFreq_ZeroTransitions(t *testing.T) {
	got := ClassifyChangeFreq(0, time.Hour)
	if got != ChangeFreqStable {
		t.Errorf("expected stable for zero transitions, got %s", got)
	}
}

func TestClassifyChangeFreq_ZeroWindow(t *testing.T) {
	got := ClassifyChangeFreq(100, 0)
	if got != ChangeFreqStable {
		t.Errorf("expected stable for zero window, got %s", got)
	}
}
