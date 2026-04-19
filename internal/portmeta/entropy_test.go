package portmeta

import "testing"

func TestEntropyLevel_String(t *testing.T) {
	cases := []struct {
		level EntropyLevel
		want  string
	}{
		{EntropyNone, "none"},
		{EntropyLow, "low"},
		{EntropyMedium, "medium"},
		{EntropyHigh, "high"},
		{EntropyLevel(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}

func TestEntropyFor_NoEvents(t *testing.T) {
	if got := EntropyFor(0, 100); got != EntropyNone {
		t.Errorf("expected EntropyNone, got %s", got)
	}
}

func TestEntropyFor_ZeroSamples(t *testing.T) {
	if got := EntropyFor(5, 0); got != EntropyNone {
		t.Errorf("expected EntropyNone, got %s", got)
	}
}

func TestEntropyFor_LowEntropy(t *testing.T) {
	// 1 event in 100 samples → very low ratio → low entropy
	got := EntropyFor(1, 100)
	if got != EntropyLow {
		t.Errorf("expected EntropyLow, got %s", got)
	}
}

func TestEntropyFor_HighEntropy(t *testing.T) {
	// 50 events in 100 samples → ratio 0.5 → maximum binary entropy
	got := EntropyFor(50, 100)
	if got != EntropyHigh {
		t.Errorf("expected EntropyHigh, got %s", got)
	}
}

func TestEntropyFor_MediumEntropy(t *testing.T) {
	// 10 events in 100 samples → ratio 0.1 → medium range
	got := EntropyFor(10, 100)
	if got != EntropyMedium {
		t.Errorf("expected EntropyMedium, got %s", got)
	}
}

func TestEntropyFor_RatioCappedAtOne(t *testing.T) {
	// More events than samples should not panic
	got := EntropyFor(200, 100)
	if got == EntropyLevel(99) {
		t.Error("unexpected unknown level")
	}
}
