package monitor_test

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func startServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestMonitor_DetectsNewPort(t *testing.T) {
	port, stop := startServer(t)
	defer stop()

	var mu sync.Mutex
	var diffs []scanner.DiffResult

	m := monitor.New(monitor.Config{
		Ports:    []int{port},
		Interval: 50 * time.Millisecond,
		Alert: func(d scanner.DiffResult) {
			mu.Lock()
			diffs = append(diffs, d)
			mu.Unlock()
		},
	})

	go m.Start()
	time.Sleep(200 * time.Millisecond)
	m.Stop()

	// No changes expected since port was open from the start and stays open.
	mu.Lock()
	defer mu.Unlock()
	if len(diffs) != 0 {
		t.Errorf("expected 0 alerts, got %d", len(diffs))
	}
}

func TestMonitor_DetectsClosedPort(t *testing.T) {
	port, stop := startServer(t)

	var mu sync.Mutex
	var diffs []scanner.DiffResult

	m := monitor.New(monitor.Config{
		Ports:    []int{port},
		Interval: 50 * time.Millisecond,
		Alert: func(d scanner.DiffResult) {
			mu.Lock()
			diffs = append(diffs, d)
			mu.Unlock()
		},
	})

	go m.Start()
	time.Sleep(80 * time.Millisecond)
	stop() // close the port mid-monitoring
	time.Sleep(150 * time.Millisecond)
	m.Stop()

	mu.Lock()
	defer mu.Unlock()
	if len(diffs) == 0 {
		t.Error("expected at least one alert after port closed, got none")
	}
}
