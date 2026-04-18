package sampler_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/sampler"
	"github.com/user/portwatch/internal/scanner"
)

func startTCP(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	return l.Addr().(*net.TCPAddr).Port
}

func TestSampler_DeliversResult(t *testing.T) {
	port := startTCP(t)
	s := scanner.NewTCPScanner("127.0.0.1", []int{port}, 200*time.Millisecond)
	smp := sampler.New(s, 50*time.Millisecond, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go smp.Run(ctx)

	select {
	case r := <-smp.Results():
		if r.Err != nil {
			t.Fatalf("unexpected error: %v", r.Err)
		}
		if len(r.Ports) == 0 {
			t.Fatal("expected at least one open port")
		}
		if r.At.IsZero() {
			t.Fatal("timestamp should be set")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for result")
	}
}

func TestSampler_StopsOnContextCancel(t *testing.T) {
	s := scanner.NewTCPScanner("127.0.0.1", []int{19999}, 100*time.Millisecond)
	smp := sampler.New(s, 500*time.Millisecond, 0)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		smp.Run(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("sampler did not stop after context cancel")
	}
}

func TestSampler_ChannelClosedAfterStop(t *testing.T) {
	s := scanner.NewTCPScanner("127.0.0.1", []int{19998}, 100*time.Millisecond)
	smp := sampler.New(s, 500*time.Millisecond, 0)

	ctx, cancel := context.WithCancel(context.Background())
	go smp.Run(ctx)
	cancel()

	timeout := time.After(time.Second)
	for {
		select {
		case _, ok := <-smp.Results():
			if !ok {
				return
			}
		case <-timeout:
			t.Fatal("results channel not closed")
		}
	}
}
