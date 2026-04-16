package watchdog_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

// TestIntegration_StallThenRecover verifies that after a stall is detected
// and beats resume, no further stall callbacks fire.
func TestIntegration_StallThenRecover(t *testing.T) {
	w := watchdog.New(60 * time.Millisecond)
	var stalls atomic.Int32
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go w.Run(ctx, func() { stalls.Add(1) })

	// Let it stall once.
	time.Sleep(150 * time.Millisecond)
	if stalls.Load() == 0 {
		t.Fatal("expected at least one stall before recovery")
	}

	// Resume beating and reset counter.
	stalls.Store(0)
	stop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				w.Beat()
			}
		}
	}()

	time.Sleep(200 * time.Millisecond)
	close(stop)

	if stalls.Load() != 0 {
		t.Fatalf("expected no stalls after recovery, got %d", stalls.Load())
	}
}
