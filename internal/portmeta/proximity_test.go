package portmeta

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makeProxPort(port int) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Proto: "tcp"}
}

func TestProximityLevel_String(t *testing.T) {
	cases := []struct {
		level ProximityLevel
		want  string
	}{
		{ProximityNone, "none"},
		{ProximityDistant, "distant"},
		{ProximityNear, "near"},
		{ProximityAdjacent, "adjacent"},
		{ProximityImmediate, "immediate"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("ProximityLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestProximityFor_NoOtherPorts(t *testing.T) {
	p := makeProxPort(8080)
	got := ProximityFor(p, []scanner.PortInfo{p})
	if got != ProximityNone {
		t.Errorf("expected ProximityNone, got %s", got)
	}
}

func TestProximityFor_ImmediateNeighbour(t *testing.T) {
	p := makeProxPort(81) // 80 is well-known HTTP
	all := []scanner.PortInfo{makeProxPort(80), p}
	got := ProximityFor(p, all)
	if got != ProximityImmediate {
		t.Errorf("expected ProximityImmediate, got %s", got)
	}
}

func TestProximityFor_NearNeighbour(t *testing.T) {
	p := makeProxPort(130) // 80 is ~50 away
	all := []scanner.PortInfo{makeProxPort(80), p}
	got := ProximityFor(p, all)
	if got != ProximityNear {
		t.Errorf("expected ProximityNear, got %s", got)
	}
}

func TestProximityAnnotator_AddsMetadata(t *testing.T) {
	p := makeProxPort(22) // SSH well-known
	target := makeProxPort(23)
	all := []scanner.PortInfo{p, target}
	annotate := NewProximityAnnotator(all)
	result := annotate([]scanner.PortInfo{target})
	if result[0].Meta == nil {
		t.Fatal("expected Meta to be set")
	}
	if _, ok := result[0].Meta[proximityKey]; !ok {
		t.Error("expected proximity key in Meta")
	}
}

func TestFilterByMinProximity_ExcludesLow(t *testing.T) {
	p := makeProxPort(9999)
	p.Meta = map[string]string{
		proximityKey:             "none",
		proximityKey + "_score": "0",
	}
	result := FilterByMinProximity([]scanner.PortInfo{p}, ProximityNear)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d ports", len(result))
	}
}

func TestFilterByMinProximity_NoMeta_PassesThrough(t *testing.T) {
	p := makeProxPort(9999)
	result := FilterByMinProximity([]scanner.PortInfo{p}, ProximityNear)
	if len(result) != 1 {
		t.Errorf("expected port to pass through, got %d", len(result))
	}
}
