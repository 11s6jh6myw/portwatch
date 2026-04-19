package portmeta

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func TestClassifyAge_Fresh(t *testing.T) {
	now := time.Now()
	if got := ClassifyAge(now.Add(-30*time.Minute), now); got != AgeFresh {
		t.Fatalf("expected fresh, got %s", got)
	}
}

func TestClassifyAge_ShortLived(t *testing.T) {
	now := time.Now()
	if got := ClassifyAge(now.Add(-2*time.Hour), now); got != AgeShortLived {
		t.Fatalf("expected short-lived, got %s", got)
	}
}

func TestClassifyAge_Mature(t *testing.T) {
	now := time.Now()
	if got := ClassifyAge(now.Add(-3*24*time.Hour), now); got != AgeMature {
		t.Fatalf("expected mature, got %s", got)
	}
}

func TestClassifyAge_Established(t *testing.T) {
	now := time.Now()
	if got := ClassifyAge(now.Add(-10*24*time.Hour), now); got != AgeEstablished {
		t.Fatalf("expected established, got %s", got)
	}
}

func TestClassifyAge_ZeroTime(t *testing.T) {
	if got := ClassifyAge(time.Time{}, time.Now()); got != AgeUnknown {
		t.Fatalf("expected unknown, got %s", got)
	}
}

func TestAgeAnnotator_SetsAgeClass(t *testing.T) {
	now := time.Now()
	a := &AgeAnnotator{now: func() time.Time { return now }}
	ports := []scanner.PortInfo{
		{Port: 80, Meta: map[string]string{"first_seen": now.Add(-2 * time.Hour).Format(time.RFC3339)}},
	}
	out := a.Annotate(ports)
	if out[0].Meta["age_class"] != "short-lived" {
		t.Fatalf("unexpected age_class: %s", out[0].Meta["age_class"])
	}
	if out[0].Meta["age_seconds"] == "" {
		t.Fatal("expected age_seconds to be set")
	}
}

func TestAgeAnnotator_NoFirstSeen(t *testing.T) {
	a := NewAgeAnnotator()
	ports := []scanner.PortInfo{{Port: 443}}
	out := a.Annotate(ports)
	if out[0].Meta["age_class"] != "unknown" {
		t.Fatalf("expected unknown, got %s", out[0].Meta["age_class"])
	}
}

func TestFilterByMaxAge_ReturnsFreshOnly(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 80, Meta: map[string]string{"age_class": "fresh"}},
		{Port: 443, Meta: map[string]string{"age_class": "established"}},
	}
	out := FilterByMaxAge(ports, AgeFresh)
	if len(out) != 1 || out[0].Port != 80 {
		t.Fatalf("unexpected result: %v", out)
	}
}
