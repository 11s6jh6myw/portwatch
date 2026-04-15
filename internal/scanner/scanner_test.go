package scanner

import (
	"net"
	"strconv"
	"testing"
)

// startTestServer opens a TCP listener on a random port and returns the port and a closer func.
func startTestServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port, _ := strconv.Atoi(ln.Addr().(*net.TCPAddr).Port.String())
	// Re-extract port correctly
	port = ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestTCPScanner_FindsOpenPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not start listener: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port

	scanner := NewTCPScanner("127.0.0.1", port, port)
	ports, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	if len(ports) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(ports))
	}

	if ports[0].Port != port {
		t.Errorf("expected port %d, got %d", port, ports[0].Port)
	}

	if ports[0].Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", ports[0].Protocol)
	}
}

func TestTCPScanner_NoOpenPorts(t *testing.T) {
	// Use a port range unlikely to be open in test environments.
	scanner := NewTCPScanner("127.0.0.1", 19999, 19999)
	ports, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(ports) != 0 {
		t.Errorf("expected 0 open ports, got %d", len(ports))
	}
}

func TestPortInfo_String(t *testing.T) {
	p := PortInfo{Protocol: "tcp", Port: 8080, Address: "127.0.0.1"}
	expected := "127.0.0.1:8080 (tcp)"
	if p.String() != expected {
		t.Errorf("expected %q, got %q", expected, p.String())
	}
}
