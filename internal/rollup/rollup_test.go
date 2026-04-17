package rollup_test

import (
	"testing"
	"time"

	"github.com/joshbeard/portwatch/internal/rollup"
	"github.com/joshbeard/portwatch/internal/scanner"
)

func port(p int) scanner.PortInfo {
	return scanner.PortInfo{Port: p, Proto: "tcp"}
}

func TestRoller_BatchesEvents(t *testing.T) {
	results := make(chan rollup.Summary, 1)
	r := rollup.New(50*time.Millisecond, func(s rollup.Summary) { results <- s })

	r.Add(rollup.Event{Port: port(80), Opened: true})
	r.Add(rollup.Event{Port: port(443), Opened: true})
	r.Add(rollup.Event{Port: port(8080), Opened: false})

	select {
	case s := <-results:
		if len(s.Opened) != 2 {
			t.Fatalf("expected 2 opened, got %d", len(s.Opened))
		}
		if len(s.Closed) != 1 {
			t.Fatalf("expected 1 closed, got %d", len(s.Closed))
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for summary")
	}
}

func TestRoller_Flush_ImmediateEmit(t *testing.T) {
	results := make(chan rollup.Summary, 1)
	r := rollup.New(10*time.Second, func(s rollup.Summary) { results <- s })

	r.Add(rollup.Event{Port: port(22), Opened: true})
	r.Flush()

	select {
	case s := <-results:
		if len(s.Opened) != 1 {
			t.Fatalf("expected 1 opened, got %d", len(s.Opened))
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for flush")
	}
}

func TestRoller_Flush_EmptyBatch(t *testing.T) {
	called := false
	r := rollup.New(10*time.Second, func(s rollup.Summary) { called = true })
	r.Flush()
	if !called {
		t.Fatal("expected onFlush to be called even for empty batch")
	}
}

func TestRoller_TimerResetAfterFlush(t *testing.T) {
	count := 0
	r := rollup.New(30*time.Millisecond, func(_ rollup.Summary) { count++ })

	r.Add(rollup.Event{Port: port(80), Opened: true})
	time.Sleep(80 * time.Millisecond)
	r.Add(rollup.Event{Port: port(443), Opened: true})
	time.Sleep(80 * time.Millisecond)

	if count != 2 {
		t.Fatalf("expected 2 flushes, got %d", count)
	}
}
