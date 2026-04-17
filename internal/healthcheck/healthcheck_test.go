package healthcheck_test

import (
	"context"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
)

func startTCP(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return ln.Addr().String()
}

func TestChecker_HealthyServer_NoFailCallback(t *testing.T) {
	addr := startTCP(t)
	var called atomic.Int32
	c := healthcheck.New(addr, 20*time.Millisecond, 200*time.Millisecond, func(_ healthcheck.Status) {
		called.Add(1)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	c.Run(ctx)
	if called.Load() != 0 {
		t.Fatalf("expected no failures, got %d", called.Load())
	}
}

func TestChecker_UnreachableServer_CallsOnFail(t *testing.T) {
	var statuses []healthcheck.Status
	c := healthcheck.New("127.0.0.1:1", 20*time.Millisecond, 50*time.Millisecond, func(s healthcheck.Status) {
		statuses = append(statuses, s)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	c.Run(ctx)
	if len(statuses) == 0 {
		t.Fatal("expected at least one failure callback")
	}
	for _, s := range statuses {
		if s.Healthy {
			t.Error("status should not be healthy")
		}
		if s.Err == nil {
			t.Error("expected non-nil error")
		}
	}
}

func TestChecker_StopsOnContextCancel(t *testing.T) {
	addr := startTCP(t)
	done := make(chan struct{})
	c := healthcheck.New(addr, 10*time.Millisecond, 100*time.Millisecond, func(_ healthcheck.Status) {})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c.Run(ctx)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Run did not stop after context cancel")
	}
}
