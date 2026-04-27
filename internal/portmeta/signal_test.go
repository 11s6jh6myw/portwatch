package portmeta

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeSignalPort(port int, firstSeen time.Time) scanner.PortInfo {
	return scanner.PortInfo{
		Port:      port,
		FirstSeen: firstSeen,
		Meta:      make(map[string]string),
	}
}

func TestSignalStrength_String(t *testing.T) {
	cases := []struct {
		s    SignalStrength
		want string
	}{
		{SignalNone, "none"},
		{SignalWeak, "weak"},
		{SignalModerate, "moderate"},
		{SignalStrong, "strong"},
		{SignalCritical, "critical"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("SignalStrength(%d).String() = %q; want %q", tc.s, got, tc.want)
		}
	}
}

func TestSignalFor_SafePort_ReturnsNoneOrWeak(t *testing.T) {
	now := time.Now()
	p := makeSignalPort(80, now.Add(-1*time.Hour))
	s := SignalFor(p, now)
	if s > SignalWeak {
		t.Errorf("expected none or weak for port 80, got %s", s)
	}
}

func TestSignalFor_HighRiskRecentPort_ReturnsCriticalOrStrong(t *testing.T) {
	now := time.Now()
	// port 23 (telnet) is high risk; first seen 1 minute ago
	p := makeSignalPort(23, now.Add(-1*time.Minute))
	s := SignalFor(p, now)
	if s < SignalStrong {
		t.Errorf("expected strong or critical for port 23 newly opened, got %s", s)
	}
}

func TestSignalAnnotator_AddsMetadata(t *testing.T) {
	now := time.Now()
	ports := []scanner.PortInfo{
		makeSignalPort(22, now.Add(-10*time.Minute)),
	}
	annotate := NewSignalAnnotator(now)
	out := annotate(ports)
	if len(out) != 1 {
		t.Fatalf("expected 1 port, got %d", len(out))
	}
	if _, ok := out[0].Meta[signalKey]; !ok {
		t.Error("expected signal metadata key to be set")
	}
}

func TestFilterByMinSignal_IncludesHighEnough(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 1, Meta: map[string]string{signalKey: "critical"}},
		{Port: 2, Meta: map[string]string{signalKey: "weak"}},
		{Port: 3, Meta: map[string]string{signalKey: "none"}},
	}
	out := FilterByMinSignal(ports, SignalModerate)
	if len(out) != 1 || out[0].Port != 1 {
		t.Errorf("expected only port 1 (critical), got %+v", out)
	}
}

func TestFilterByMinSignal_NoMeta_PassesThrough(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 9, Meta: nil},
	}
	out := FilterByMinSignal(ports, SignalStrong)
	if len(out) != 1 {
		t.Errorf("expected port with no meta to pass through, got %d ports", len(out))
	}
}
