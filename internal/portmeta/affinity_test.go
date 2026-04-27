package portmeta

import (
	"testing"

	"github.com/iamcalledrob/portwatch/internal/scanner"
)

func makeAffinityPort(port int) scanner.PortInfo {
	return scanner.PortInfo{Port: port}
}

func TestAffinityLevel_String(t *testing.T) {
	cases := []struct {
		level AffinityLevel
		want  string
	}{
		{AffinityNone, "none"},
		{AffinityWeak, "weak"},
		{AffinityMedium, "medium"},
		{AffinityStrong, "strong"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("AffinityLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestAffinityFor_StrongCanonicalPort(t *testing.T) {
	level, family := AffinityFor(makeAffinityPort(443))
	if level != AffinityStrong {
		t.Errorf("port 443: got level %s, want strong", level)
	}
	if family != "web" {
		t.Errorf("port 443: got family %q, want web", family)
	}
}

func TestAffinityFor_WeakAlternatePort(t *testing.T) {
	level, family := AffinityFor(makeAffinityPort(8080))
	if level != AffinityWeak {
		t.Errorf("port 8080: got level %s, want weak", level)
	}
	if family != "web" {
		t.Errorf("port 8080: got family %q, want web", family)
	}
}

func TestAffinityFor_UnknownPort(t *testing.T) {
	level, family := AffinityFor(makeAffinityPort(9999))
	if level != AffinityNone {
		t.Errorf("port 9999: got level %s, want none", level)
	}
	if family != "" {
		t.Errorf("port 9999: expected empty family, got %q", family)
	}
}

func TestAffinityAnnotator_AddsMetadata(t *testing.T) {
	annotate := NewAffinityAnnotator()
	ports := []scanner.PortInfo{makeAffinityPort(22), makeAffinityPort(9999)}
	out := annotate(ports)

	if out[0].Meta["affinity"] != "strong" {
		t.Errorf("port 22: affinity = %q, want strong", out[0].Meta["affinity"])
	}
	if out[0].Meta["affinity.family"] != "remote" {
		t.Errorf("port 22: family = %q, want remote", out[0].Meta["affinity.family"])
	}
	if out[1].Meta["affinity"] != "none" {
		t.Errorf("port 9999: affinity = %q, want none", out[1].Meta["affinity"])
	}
	if _, ok := out[1].Meta["affinity.family"]; ok {
		t.Error("port 9999: unexpected affinity.family key")
	}
}

func TestFilterByMinAffinity_ExcludesBelow(t *testing.T) {
	annotate := NewAffinityAnnotator()
	ports := annotate([]scanner.PortInfo{
		makeAffinityPort(80),   // strong
		makeAffinityPort(8080), // weak
		makeAffinityPort(9999), // none
	})

	result := FilterByMinAffinity(ports, AffinityMedium)
	if len(result) != 1 {
		t.Fatalf("expected 1 port, got %d", len(result))
	}
	if result[0].Port != 80 {
		t.Errorf("expected port 80, got %d", result[0].Port)
	}
}
