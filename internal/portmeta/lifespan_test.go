package portmeta

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func TestLifespanLevel_String(t *testing.T) {
	cases := []struct {
		level LifespanLevel
		want  string
	}{
		{LifespanUnknown, "unknown"},
		{LifespanEphemeral, "ephemeral"},
		{LifespanShort, "short"},
		{LifespanMedium, "medium"},
		{LifespanLong, "long"},
		{LifespanPermanent, "permanent"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("LifespanLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestLifespanFor_ZeroTime(t *testing.T) {
	if got := LifespanFor(time.Time{}); got != LifespanUnknown {
		t.Errorf("expected LifespanUnknown for zero time, got %v", got)
	}
}

func TestLifespanFor_Ephemeral(t *testing.T) {
	firstSeen := time.Now().Add(-10 * time.Second)
	if got := LifespanFor(firstSeen); got != LifespanEphemeral {
		t.Errorf("expected LifespanEphemeral, got %v", got)
	}
}

func TestLifespanFor_Short(t *testing.T) {
	firstSeen := time.Now().Add(-30 * time.Minute)
	if got := LifespanFor(firstSeen); got != LifespanShort {
		t.Errorf("expected LifespanShort, got %v", got)
	}
}

func TestLifespanFor_Long(t *testing.T) {
	firstSeen := time.Now().Add(-48 * time.Hour)
	if got := LifespanFor(firstSeen); got != LifespanLong {
		t.Errorf("expected LifespanLong, got %v", got)
	}
}

func TestLifespanFor_Permanent(t *testing.T) {
	firstSeen := time.Now().Add(-10 * 24 * time.Hour)
	if got := LifespanFor(firstSeen); got != LifespanPermanent {
		t.Errorf("expected LifespanPermanent, got %v", got)
	}
}

func TestIsLongLived_True(t *testing.T) {
	if !IsLongLived(time.Now().Add(-48 * time.Hour)) {
		t.Error("expected IsLongLived=true for 48h old port")
	}
}

func TestIsLongLived_False(t *testing.T) {
	if IsLongLived(time.Now().Add(-30 * time.Minute)) {
		t.Error("expected IsLongLived=false for 30m old port")
	}
}

func TestLifespanAnnotator_AddsMetadata(t *testing.T) {
	firstSeen := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	ports := []scanner.PortInfo{
		{Port: 80, Meta: map[string]string{metaFirstSeen: firstSeen}},
	}
	annotate := NewLifespanAnnotator()
	out := annotate(ports)
	if got := out[0].Meta[metaLifespan]; got != "medium" {
		t.Errorf("expected lifespan=medium, got %q", got)
	}
}

func TestFilterByMinLifespan_IncludesLongEnough(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 22, Meta: map[string]string{metaLifespan: "permanent"}},
		{Port: 8080, Meta: map[string]string{metaLifespan: "ephemeral"}},
	}
	out := FilterByMinLifespan(ports, LifespanLong)
	if len(out) != 1 || out[0].Port != 22 {
		t.Errorf("expected only port 22, got %v", out)
	}
}

func TestFilterByMinLifespan_NoMeta_PassesThrough(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 443, Meta: map[string]string{}},
	}
	out := FilterByMinLifespan(ports, LifespanLong)
	if len(out) != 1 {
		t.Errorf("expected port without meta to pass through, got %v", out)
	}
}
