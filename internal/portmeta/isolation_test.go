package portmeta

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeIsolationPort(port uint16) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Proto: "tcp"}
}

func TestIsolationLevel_String(t *testing.T) {
	cases := []struct {
		level IsolationLevel
		want  string
	}{
		{IsolationNone, "none"},
		{IsolationLow, "low"},
		{IsolationMedium, "medium"},
		{IsolationHigh, "high"},
		{IsolationLevel(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("IsolationLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestIsolationFor_NoPeers_ReturnsHigh(t *testing.T) {
	p := makeIsolationPort(8080)
	level := IsolationFor(p, nil, time.Now())
	if level != IsolationHigh {
		t.Errorf("expected IsolationHigh, got %s", level)
	}
}

func TestIsolationFor_ManyNearbyPeers_ReturnsNone(t *testing.T) {
	p := makeIsolationPort(8080)
	peers := []scanner.PortInfo{
		makeIsolationPort(8000),
		makeIsolationPort(8010),
		makeIsolationPort(8020),
		makeIsolationPort(8030),
		makeIsolationPort(8040),
		makeIsolationPort(8050),
	}
	level := IsolationFor(p, peers, time.Now())
	if level != IsolationNone {
		t.Errorf("expected IsolationNone, got %s", level)
	}
}

func TestIsolationFor_OneNearbyPeer_ReturnsMedium(t *testing.T) {
	p := makeIsolationPort(9000)
	peers := []scanner.PortInfo{
		makeIsolationPort(9050),
	}
	level := IsolationFor(p, peers, time.Now())
	if level != IsolationMedium {
		t.Errorf("expected IsolationMedium, got %s", level)
	}
}

func TestIsolationFor_IgnoresSelf(t *testing.T) {
	p := makeIsolationPort(7777)
	peers := []scanner.PortInfo{makeIsolationPort(7777)}
	level := IsolationFor(p, peers, time.Now())
	if level != IsolationHigh {
		t.Errorf("expected IsolationHigh when only self in peers, got %s", level)
	}
}

func TestIsolationFor_DistantPeersOnly_ReturnsHigh(t *testing.T) {
	p := makeIsolationPort(80)
	peers := []scanner.PortInfo{
		makeIsolationPort(9000),
		makeIsolationPort(9100),
	}
	level := IsolationFor(p, peers, time.Now())
	if level != IsolationHigh {
		t.Errorf("expected IsolationHigh for distant peers, got %s", level)
	}
}
