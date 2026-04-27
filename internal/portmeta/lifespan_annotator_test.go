package portmeta

import (
	"fmt"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeLifespanPort(port int, firstSeenAgo time.Duration) scanner.PortInfo {
	meta := map[string]string{}
	if firstSeenAgo > 0 {
		meta[metaFirstSeen] = time.Now().Add(-firstSeenAgo).Format(time.RFC3339)
	}
	return scanner.PortInfo{Port: port, Meta: meta}
}

func TestLifespanAnnotator_NilMetaInitialised(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 9090},
	}
	annotate := NewLifespanAnnotator()
	out := annotate(ports)
	if out[0].Meta == nil {
		t.Fatal("expected Meta to be initialised")
	}
	if _, ok := out[0].Meta[metaLifespan]; !ok {
		t.Error("expected lifespan key to be set")
	}
}

func TestLifespanAnnotator_MediumPort(t *testing.T) {
	p := makeLifespanPort(443, 6*time.Hour)
	annotate := NewLifespanAnnotator()
	out := annotate([]scanner.PortInfo{p})
	if got := out[0].Meta[metaLifespan]; got != "medium" {
		t.Errorf("expected medium, got %q", got)
	}
}

func TestFilterByMinLifespan_ExcludesShort(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 80, Meta: map[string]string{metaLifespan: "short"}},
		{Port: 22, Meta: map[string]string{metaLifespan: "long"}},
	}
	out := FilterByMinLifespan(ports, LifespanMedium)
	if len(out) != 1 || out[0].Port != 22 {
		t.Errorf("expected only port 22, got %v", out)
	}
}

func TestParseLifespan_RoundTrip(t *testing.T) {
	levels := []LifespanLevel{
		LifespanEphemeral,
		LifespanShort,
		LifespanMedium,
		LifespanLong,
		LifespanPermanent,
	}
	for _, l := range levels {
		if got := parseLifespan(l.String()); got != l {
			t.Errorf("parseLifespan(%q) = %v, want %v", l.String(), got, l)
		}
	}
}

func TestFirstSeenFromMeta_Valid(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	meta := map[string]string{metaFirstSeen: now.Format(time.RFC3339)}
	got := firstSeenFromMeta(meta)
	if !got.Equal(now) {
		t.Errorf("expected %v, got %v", now, got)
	}
}

func TestFirstSeenFromMeta_Missing(t *testing.T) {
	got := firstSeenFromMeta(map[string]string{})
	if !got.IsZero() {
		t.Errorf("expected zero time for missing key, got %v", got)
	}
}

func TestLifespanAnnotator_AllLevels(t *testing.T) {
	cases := []struct {
		ago  time.Duration
		want string
	}{
		{10 * time.Second, "ephemeral"},
		{30 * time.Minute, "short"},
		{6 * time.Hour, "medium"},
		{48 * time.Hour, "long"},
		{10 * 24 * time.Hour, "permanent"},
	}
	annotate := NewLifespanAnnotator()
	for _, tc := range cases {
		p := makeLifespanPort(80, tc.ago)
		out := annotate([]scanner.PortInfo{p})
		if got := out[0].Meta[metaLifespan]; got != tc.want {
			t.Errorf("ago=%v: expected %q, got %q", tc.ago, tc.want, got)
		}
	}
	_ = fmt.Sprintf // suppress import
}
