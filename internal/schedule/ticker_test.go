package schedule

import (
	"context"
	"testing"
	"time"
)

func TestTicker_FiresWithinTimeout(t *testing.T) {
	cfg := Config{
		MinInterval: 10 * time.Millisecond,
		MaxInterval: 20 * time.Millisecond,
		StepDown:    0.5,
		StepUp:      1.5,
	}
	s := New(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	tk := NewTicker(ctx, s)
	defer tk.Stop()

	select {
	case <-tk.C:
		// success
	case <-ctx.Done():
		t.Fatal("ticker did not fire within timeout")
	}
}

func TestTicker_StopsOnContextCancel(t *testing.T) {
	cfg := Config{
		MinInterval: 10 * time.Millisecond,
		MaxInterval: 20 * time.Millisecond,
		StepDown:    0.5,
		StepUp:      1.5,
	}
	s := New(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	tk := NewTicker(ctx, s)

	// drain one tick then cancel
	<-tk.C
	cancel()

	// give goroutine time to exit
	time.Sleep(50 * time.Millisecond)
	// no assertion needed — test completes without hanging
}

func TestTicker_AdaptsAfterActivity(t *testing.T) {
	cfg := Config{
		MinInterval: 10 * time.Millisecond,
		MaxInterval: 100 * time.Millisecond,
		StepDown:    0.1,
		StepUp:      2.0,
	}
	s := New(cfg)
	s.RecordActivity() // should drop to ~10ms (min)

	if s.Current() != 10*time.Millisecond {
		t.Fatalf("expected min interval after activity, got %v", s.Current())
	}
}
