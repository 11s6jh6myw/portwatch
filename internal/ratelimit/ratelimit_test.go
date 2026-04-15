package ratelimit

import (
	"testing"
	"time"
)

func TestAllow_FirstEventAlwaysAllowed(t *testing.T) {
	l := New(5 * time.Second)
	if !l.Allow("tcp", 8080, "opened") {
		t.Fatal("expected first event to be allowed")
	}
}

func TestAllow_DuplicateWithinCooldownSuppressed(t *testing.T) {
	now := time.Now()
	l := New(10 * time.Second)
	l.now = func() time.Time { return now }

	if !l.Allow("tcp", 8080, "opened") {
		t.Fatal("expected first event to be allowed")
	}
	if l.Allow("tcp", 8080, "opened") {
		t.Fatal("expected duplicate within cooldown to be suppressed")
	}
}

func TestAllow_DuplicateAfterCooldownAllowed(t *testing.T) {
	now := time.Now()
	l := New(5 * time.Second)
	l.now = func() time.Time { return now }

	l.Allow("tcp", 8080, "opened")

	l.now = func() time.Time { return now.Add(6 * time.Second) }
	if !l.Allow("tcp", 8080, "opened") {
		t.Fatal("expected event after cooldown to be allowed")
	}
}

func TestAllow_DifferentEventsAreIndependent(t *testing.T) {
	l := New(10 * time.Second)
	l.Allow("tcp", 8080, "opened")
	if !l.Allow("tcp", 8080, "closed") {
		t.Fatal("expected different event type to be allowed")
	}
}

func TestAllow_DifferentPortsAreIndependent(t *testing.T) {
	l := New(10 * time.Second)
	l.Allow("tcp", 8080, "opened")
	if !l.Allow("tcp", 9090, "opened") {
		t.Fatal("expected different port to be allowed")
	}
}

func TestFlush_ResetsState(t *testing.T) {
	now := time.Now()
	l := New(10 * time.Second)
	l.now = func() time.Time { return now }

	l.Allow("tcp", 8080, "opened")
	l.Flush()

	if !l.Allow("tcp", 8080, "opened") {
		t.Fatal("expected event to be allowed after flush")
	}
}

func TestExpire_RemovesStaleEntries(t *testing.T) {
	now := time.Now()
	l := New(5 * time.Second)
	l.now = func() time.Time { return now }

	l.Allow("tcp", 8080, "opened")

	l.now = func() time.Time { return now.Add(6 * time.Second) }
	l.Expire()

	if len(l.seen) != 0 {
		t.Fatalf("expected stale entries to be removed, got %d", len(l.seen))
	}
}
