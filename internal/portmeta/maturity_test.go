package portmeta

import (
	"testing"
	"time"
)

func TestMaturityLevel_String(t *testing.T) {
	cases := []struct {
		level MaturityLevel
		want  string
	}{
		{MaturityUnknown, "unknown"},
		{MaturityEmerging, "emerging"},
		{MaturityDeveloping, "developing"},
		{MaturityEstablished, "established"},
		{MaturityMature, "mature"},
	}
	for _, c := range cases {
		if got := c.level.String(); got != c.want {
			t.Errorf("String() = %q, want %q", got, c.want)
		}
	}
}

func TestMaturityFor_ZeroTime(t *testing.T) {
	if got := MaturityFor(time.Time{}, 0); got != MaturityUnknown {
		t.Errorf("expected unknown, got %s", got)
	}
}

func TestMaturityFor_Emerging(t *testing.T) {
	first := time.Now().Add(-1 * time.Hour)
	if got := MaturityFor(first, 2); got != MaturityEmerging {
		t.Errorf("expected emerging, got %s", got)
	}
}

func TestMaturityFor_Developing(t *testing.T) {
	first := time.Now().Add(-3 * 24 * time.Hour)
	if got := MaturityFor(first, 10); got != MaturityDeveloping {
		t.Errorf("expected developing, got %s", got)
	}
}

func TestMaturityFor_Established(t *testing.T) {
	first := time.Now().Add(-15 * 24 * time.Hour)
	if got := MaturityFor(first, 50); got != MaturityEstablished {
		t.Errorf("expected established, got %s", got)
	}
}

func TestMaturityFor_Mature(t *testing.T) {
	first := time.Now().Add(-60 * 24 * time.Hour)
	if got := MaturityFor(first, 200); got != MaturityMature {
		t.Errorf("expected mature, got %s", got)
	}
}
