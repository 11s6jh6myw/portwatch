package debounce_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/debounce"
)

func TestAllow_FirstEventAlwaysAllowed(t *testing.T) {
	d := debounce.New(1 * time.Second)
	if !d.Allow(8080, "opened") {
		t.Fatal("expected first event to be allowed")
	}
}

func TestAllow_DuplicateWithinWindowSuppressed(t *testing.T) {
	d := debounce.New(5 * time.Second)
	d.Allow(8080, "opened")
	if d.Allow(8080, "opened") {
		t.Fatal("expected duplicate within window to be suppressed")
	}
}

func TestAllow_DuplicateAfterWindowAllowed(t *testing.T) {
	now := time.Now()
	d := debounce.New(100 * time.Millisecond)

	// Override clock to control time.
	calls := 0
	d = debounce.New(100 * time.Millisecond)
	_ = calls
	_ = now

	d.Allow(9090, "closed")
	time.Sleep(150 * time.Millisecond)

	if !d.Allow(9090, "closed") {
		t.Fatal("expected event after window expiry to be allowed")
	}
}

func TestAllow_DifferentEventsAreIndependent(t *testing.T) {
	d := debounce.New(5 * time.Second)
	if !d.Allow(8080, "opened") {
		t.Fatal("expected opened to be allowed")
	}
	if !d.Allow(8080, "closed") {
		t.Fatal("expected closed to be allowed independently")
	}
}

func TestAllow_DifferentPortsAreIndependent(t *testing.T) {
	d := debounce.New(5 * time.Second)
	if !d.Allow(8080, "opened") {
		t.Fatal("expected port 8080 to be allowed")
	}
	if !d.Allow(9090, "opened") {
		t.Fatal("expected port 9090 to be allowed independently")
	}
}

func TestPending_ReflectsActiveEntries(t *testing.T) {
	d := debounce.New(5 * time.Second)
	if d.Pending() != 0 {
		t.Fatalf("expected 0 pending, got %d", d.Pending())
	}
	d.Allow(8080, "opened")
	d.Allow(9090, "closed")
	if got := d.Pending(); got != 2 {
		t.Fatalf("expected 2 pending, got %d", got)
	}
}

func TestPending_EvictsExpiredEntries(t *testing.T) {
	d := debounce.New(50 * time.Millisecond)
	d.Allow(1234, "opened")
	time.Sleep(80 * time.Millisecond)
	if got := d.Pending(); got != 0 {
		t.Fatalf("expected 0 pending after expiry, got %d", got)
	}
}
