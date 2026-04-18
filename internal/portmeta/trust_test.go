package portmeta

import (
	"testing"

	"github.com/joshbeard/portwatch/internal/scanner"
)

func TestTrustFor_HighTrust(t *testing.T) {
	if got := TrustFor(443); got != TrustHigh {
		t.Fatalf("expected TrustHigh, got %s", got)
	}
}

func TestTrustFor_LowTrust(t *testing.T) {
	if got := TrustFor(23); got != TrustLow {
		t.Fatalf("expected TrustLow, got %s", got)
	}
}

func TestTrustFor_Unknown(t *testing.T) {
	if got := TrustFor(9999); got != TrustUnknown {
		t.Fatalf("expected TrustUnknown, got %s", got)
	}
}

func TestIsTrusted_True(t *testing.T) {
	if !IsTrusted(80) {
		t.Fatal("expected port 80 to be trusted")
	}
}

func TestIsTrusted_False(t *testing.T) {
	if IsTrusted(4444) {
		t.Fatal("expected port 4444 to be untrusted")
	}
}

func TestTrustLevel_String(t *testing.T) {
	cases := []struct {
		level TrustLevel
		want  string
	}{
		{TrustHigh, "high"},
		{TrustMedium, "medium"},
		{TrustLow, "low"},
		{TrustUnknown, "unknown"},
	}
	for _, c := range cases {
		if got := c.level.String(); got != c.want {
			t.Errorf("String() = %q, want %q", got, c.want)
		}
	}
}

func TestTrustAnnotator_SetsMetaField(t *testing.T) {
	ports := []scanner.PortInfo{{Port: 443}, {Port: 23}}
	a := NewTrustAnnotator()
	out := a.Annotate(ports)
	if out[0].Meta[trustKey] != "high" {
		t.Errorf("expected high, got %s", out[0].Meta[trustKey])
	}
	if out[1].Meta[trustKey] != "low" {
		t.Errorf("expected low, got %s", out[1].Meta[trustKey])
	}
}

func TestFilterByMinTrust_FiltersCorrectly(t *testing.T) {
	ports := []scanner.PortInfo{{Port: 443}, {Port: 23}, {Port: 9999}}
	out := FilterByMinTrust(ports, TrustMedium)
	if len(out) != 1 || out[0].Port != 443 {
		t.Errorf("unexpected result: %+v", out)
	}
}
