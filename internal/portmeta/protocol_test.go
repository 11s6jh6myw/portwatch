package portmeta_test

import (
	"testing"

	"github.com/user/portwatch/internal/portmeta"
)

func TestProtocolFor_TCPPort(t *testing.T) {
	if got := portmeta.ProtocolFor(22); got != portmeta.ProtocolTCP {
		t.Fatalf("expected tcp, got %s", got)
	}
}

func TestProtocolFor_UDPPort(t *testing.T) {
	if got := portmeta.ProtocolFor(53); got != portmeta.ProtocolUDP {
		t.Fatalf("expected udp, got %s", got)
	}
}

func TestProtocolFor_UnknownPort(t *testing.T) {
	if got := portmeta.ProtocolFor(9999); got != portmeta.ProtocolUnknown {
		t.Fatalf("expected unknown, got %s", got)
	}
}

func TestProtocol_String(t *testing.T) {
	tests := []struct {
		p    portmeta.Protocol
		want string
	}{
		{portmeta.ProtocolTCP, "tcp"},
		{portmeta.ProtocolUDP, "udp"},
		{portmeta.ProtocolUnknown, "unknown"},
	}
	for _, tt := range tests {
		if got := tt.p.String(); got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}

func TestIsTCP(t *testing.T) {
	if !portmeta.IsTCP(443) {
		t.Fatal("expected port 443 to be TCP")
	}
	if portmeta.IsTCP(53) {
		t.Fatal("expected port 53 not to be TCP")
	}
}

func TestIsUDP(t *testing.T) {
	if !portmeta.IsUDP(161) {
		t.Fatal("expected port 161 to be UDP")
	}
	if portmeta.IsUDP(80) {
		t.Fatal("expected port 80 not to be UDP")
	}
}
