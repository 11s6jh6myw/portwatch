package digest_test

import (
	"testing"

	"github.com/user/portwatch/internal/digest"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(port int, proto string) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Protocol: proto}
}

func TestCompute_EmptySlice(t *testing.T) {
	d := digest.Compute(nil)
	if d == "" {
		t.Fatal("expected non-empty digest for empty slice")
	}
}

func TestCompute_Deterministic(t *testing.T) {
	ports := []scanner.PortInfo{makePort(80, "tcp"), makePort(443, "tcp")}
	a := digest.Compute(ports)
	b := digest.Compute(ports)
	if !digest.Equal(a, b) {
		t.Fatalf("expected equal digests, got %s vs %s", a, b)
	}
}

func TestCompute_OrderIndependent(t *testing.T) {
	a := digest.Compute([]scanner.PortInfo{makePort(80, "tcp"), makePort(443, "tcp")})
	b := digest.Compute([]scanner.PortInfo{makePort(443, "tcp"), makePort(80, "tcp")})
	if !digest.Equal(a, b) {
		t.Fatalf("order should not affect digest: %s vs %s", a, b)
	}
}

func TestCompute_DifferentPorts(t *testing.T) {
	a := digest.Compute([]scanner.PortInfo{makePort(80, "tcp")})
	b := digest.Compute([]scanner.PortInfo{makePort(8080, "tcp")})
	if digest.Equal(a, b) {
		t.Fatal("expected different digests for different ports")
	}
}

func TestCompute_DifferentProtocols(t *testing.T) {
	a := digest.Compute([]scanner.PortInfo{makePort(53, "tcp")})
	b := digest.Compute([]scanner.PortInfo{makePort(53, "udp")})
	if digest.Equal(a, b) {
		t.Fatal("expected different digests for different protocols")
	}
}

func TestEqual_Reflexive(t *testing.T) {
	ports := []scanner.PortInfo{makePort(22, "tcp")}
	d := digest.Compute(ports)
	if !digest.Equal(d, d) {
		t.Fatal("digest should equal itself")
	}
}
