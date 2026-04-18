package portmeta_test

import (
	"strings"
	"testing"

	"github.com/example/portwatch/internal/portmeta"
)

func TestLookup_KnownPort(t *testing.T) {
	m, ok := portmeta.Lookup(80)
	if !ok {
		t.Fatal("expected port 80 to be known")
	}
	if m.Service != "HTTP" {
		t.Errorf("expected HTTP, got %s", m.Service)
	}
	if m.Protocol != "tcp" {
		t.Errorf("expected tcp, got %s", m.Protocol)
	}
}

func TestLookup_UnknownPort(t *testing.T) {
	_, ok := portmeta.Lookup(9999)
	if ok {
		t.Fatal("expected port 9999 to be unknown")
	}
}

func TestIsRisky_RiskyPort(t *testing.T) {
	riskyPorts := []uint16{21, 23, 445, 3306, 3389, 6379, 27017}
	for _, p := range riskyPorts {
		if !portmeta.IsRisky(p) {
			t.Errorf("expected port %d to be risky", p)
		}
	}
}

func TestIsRisky_SafePort(t *testing.T) {
	safePorts := []uint16{22, 80, 443, 5432}
	for _, p := range safePorts {
		if portmeta.IsRisky(p) {
			t.Errorf("expected port %d not to be risky", p)
		}
	}
}

func TestIsRisky_UnknownPort(t *testing.T) {
	if portmeta.IsRisky(9999) {
		t.Error("unknown port should not be risky")
	}
}

func TestMeta_String(t *testing.T) {
	m, _ := portmeta.Lookup(443)
	s := m.String()
	if !strings.Contains(s, "443") {
		t.Errorf("expected port number in string, got %s", s)
	}
	if !strings.Contains(s, "HTTPS") {
		t.Errorf("expected service name in string, got %s", s)
	}
}
