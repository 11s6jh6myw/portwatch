package fingerprint_test

import (
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(port int, proto, state string) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Protocol: proto, State: state}
}

func TestCompute_EmptySlice(t *testing.T) {
	f := fingerprint.Compute(nil)
	if f == "" {
		t.Fatal("expected non-empty fingerprint for empty input")
	}
}

func TestCompute_Deterministic(t *testing.T) {
	ports := []scanner.PortInfo{makePort(80, "tcp", "open"), makePort(443, "tcp", "open")}
	if fingerprint.Compute(ports) != fingerprint.Compute(ports) {
		t.Fatal("fingerprint is not deterministic")
	}
}

func TestCompute_OrderIndependent(t *testing.T) {
	a := []scanner.PortInfo{makePort(80, "tcp", "open"), makePort(443, "tcp", "open")}
	b := []scanner.PortInfo{makePort(443, "tcp", "open"), makePort(80, "tcp", "open")}
	if fingerprint.Compute(a) != fingerprint.Compute(b) {
		t.Fatal("fingerprint should be order-independent")
	}
}

func TestCompute_DifferentPorts(t *testing.T) {
	a := fingerprint.Compute([]scanner.PortInfo{makePort(80, "tcp", "open")})
	b := fingerprint.Compute([]scanner.PortInfo{makePort(8080, "tcp", "open")})
	if fingerprint.Equal(a, b) {
		t.Fatal("different ports should produce different fingerprints")
	}
}

func TestCompute_DifferentState(t *testing.T) {
	a := fingerprint.Compute([]scanner.PortInfo{makePort(80, "tcp", "open")})
	b := fingerprint.Compute([]scanner.PortInfo{makePort(80, "tcp", "closed")})
	if fingerprint.Equal(a, b) {
		t.Fatal("different states should produce different fingerprints")
	}
}

func TestEqual(t *testing.T) {
	ports := []scanner.PortInfo{makePort(22, "tcp", "open")}
	f := fingerprint.Compute(ports)
	if !fingerprint.Equal(f, f) {
		t.Fatal("Equal should return true for same fingerprint")
	}
}
