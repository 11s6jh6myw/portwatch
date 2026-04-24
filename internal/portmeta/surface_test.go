package portmeta

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestAttackSurface_String(t *testing.T) {
	cases := []struct {
		level AttackSurface
		want  string
	}{
		{SurfaceNone, "none"},
		{SurfaceMinimal, "minimal"},
		{SurfaceModerate, "moderate"},
		{SurfaceSignificant, "significant"},
		{SurfaceCritical, "critical"},
		{AttackSurface(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("AttackSurface(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestAttackSurfaceFor_TelnetIsHigh(t *testing.T) {
	// port 23 (telnet) is high risk, public, and highly exposed
	got := AttackSurfaceFor(23)
	if got < SurfaceSignificant {
		t.Errorf("AttackSurfaceFor(23) = %s, want >= significant", got)
	}
}

func TestAttackSurfaceFor_LoopbackLike_LowSurface(t *testing.T) {
	// port 1 is unknown, internal, low risk — expect minimal or none
	got := AttackSurfaceFor(1)
	if got > SurfaceModerate {
		t.Errorf("AttackSurfaceFor(1) = %s, want <= moderate", got)
	}
}

func TestSurfaceAnnotator_AddsMetadata(t *testing.T) {
	annotate := NewSurfaceAnnotator()
	ports := []scanner.PortInfo{{Port: 23}}
	out := annotate(ports)
	if len(out) != 1 {
		t.Fatalf("expected 1 port, got %d", len(out))
	}
	val, ok := out[0].Meta[surfaceKey]
	if !ok {
		t.Fatal("expected attack_surface key in Meta")
	}
	if val == "" {
		t.Error("expected non-empty attack_surface value")
	}
}

func TestFilterByMinSurface_IncludesHighEnough(t *testing.T) {
	annotate := NewSurfaceAnnotator()
	ports := annotate([]scanner.PortInfo{
		{Port: 23},  // high surface
		{Port: 1},   // low surface
	})
	result := FilterByMinSurface(ports, SurfaceSignificant)
	for _, p := range result {
		level := parseSurface(p.Meta[surfaceKey])
		if level < SurfaceSignificant {
			t.Errorf("port %d has surface %s, below minimum significant", p.Port, level)
		}
	}
}

func TestFilterByMinSurface_NoMeta_PassesThrough(t *testing.T) {
	ports := []scanner.PortInfo{{Port: 9999}}
	out := FilterByMinSurface(ports, SurfaceMinimal)
	if len(out) != 1 {
		t.Errorf("expected port without meta to pass through, got %d results", len(out))
	}
}
