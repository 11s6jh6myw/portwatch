package watchdog_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

func TestBeat_UpdatesLastBeat(t *testing.T) {
	w := watchdog.New(time.Second)
	before := w.LastBeat()
	time.Sleep(5 * time.Millisecond)
	w.Beat()
	if !w.LastBeat().After(before) {
		t.Fatal("expected LastBeat to advance after Beat()")
	}
}

func TestRun_CallsOnStallWhenNoBeats(t *testing.T) {
	w := watchdog.New(50 * time.Millisecond)
	var calls atomic.Int32
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

.Run(ctx, func() { calls.Add(1) })
	time.Sleep(300 * time.Millisecond)

	if calls.Load() ==	t.Fatal("expected onStall to be called at least once")
	}
}

func TestRun_NoStallWhenBeating(t *testing.T) {
	w := watchdog.New(100 * time.Millisecond)
	var calls atomic.Int32
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()

	go w.Run(ctx, func() { calls.Add(1) })

	// send beats every 30 ms — well within the 100 ms timeout
	ticker := time.NewTicker(30 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			if calls.Load() != 0 {
				t.Fatalf("expected no stall calls, got %d", calls.Load())
			return
		case <-ticker.C:
			w.Beat()
		}
	}
}

func TestRun_StopsOnContextCancel(t *testing.T).New(50 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		w.Run(ctx, func() {})
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Run did not return after context cancel")
	}
}
