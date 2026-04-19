package portmeta

import (
	"testing"
	"time"

	"github.com/iamcathal/portwatch/internal/scanner"
)

func TestClassifyFreshness_New(t *testing.T) {
	if got := ClassifyFreshness(time.Now().Add(-1 * time.Hour)); got != FreshnessNew {
		t.Fatalf("expected new, got %s", got)
	}
}

func TestClassifyFreshness_Recent(t *testing.T) {
	if got := ClassifyFreshness(time.Now().Add(-3 * 24 * time.Hour)); got != FreshnessRecent {
		t.Fatalf("expected recent, got %s", got)
	}
}

func TestClassifyFreshness_Mature(t *testing.T) {
	if got := ClassifyFreshness(time.Now().Add(-15 * 24 * time.Hour)); got != FreshnessMature {
		t.Fatalf("expected mature, got %s", got)
	}
}

func TestClassifyFreshness_Stale(t *testing.T) {
	if got := ClassifyFreshness(time.Now().Add(-60 * 24 * time.Hour)); got != FreshnessStale {
		t.Fatalf("expected stale, got %s", got)
	}
}

func TestClassifyFreshness_ZeroTime(t *testing.T) {
	if got := ClassifyFreshness(time.Time{}); got != FreshnessUnknown {
		t.Fatalf("expected unknown, got %s", got)
	}
}

func TestFreshnessLevel_String(t *testing.T) {
	cases := map[FreshnessLevel]string{
		FreshnessNew:     "new",
		FreshnessRecent:  "recent",
		FreshnessMature:  "mature",
		FreshnessStale:   "stale",
		FreshnessUnknown: "unknown",
	}
	for lvl, want := range cases {
		if got := lvl.String(); got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	}
}

func TestFreshnessAnnotator_AddsMetadata(t *testing.T) {
	firstSeen := map[int]time.Time{
		80: time.Now().Add(-1 * time.Hour),
	}
	a := NewFreshnessAnnotator(firstSeen)
	ports := []scanner.PortInfo{{Port: 80}}
	out := a.Annotate(ports)
	if out[0].Meta[freshnessKey] != "new" {
		t.Fatalf("expected new, got %s", out[0].Meta[freshnessKey])
	}
}

func TestFilterByMinFreshness_ExcludesStale(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 80, Meta: map[string]string{freshnessKey: "new"}},
		{Port: 443, Meta: map[string]string{freshnessKey: "stale"}},
	}
	out := FilterByMinFreshness(ports, FreshnessRecent)
	if len(out) != 1 || out[0].Port != 80 {
		t.Fatalf("unexpected result: %v", out)
	}
}
