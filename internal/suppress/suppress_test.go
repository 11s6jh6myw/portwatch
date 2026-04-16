package suppress

import (
	"testing"
	"time"
)

func TestAllow_FirstEventAlwaysAllowed(t *testing.T) {
	f := New(time.Minute, 2)
	if !f.Allow(8080, Opened) {
		t.Fatal("expected first event to be allowed")
	}
}

func TestAllow_BelowThresholdAllowed(t *testing.T) {
	f := New(time.Minute, 3)
	port := uint16(9000)
	for i := 0; i < 3; i++ {
		et := Opened
		if i%2 == 1 {
			et = Closed
		}
		if !f.Allow(port, et) {
			t.Fatalf("expected event %d to be allowed", i)
		}
	}
}

func TestAllow_ExceedsThresholdSuppressed(t *testing.T) {
	f := New(time.Minute, 2)
	port := uint16(443)
	allowed := 0
	for i := 0; i < 6; i++ {
		et := Opened
		if i%2 == 1 {
			et = Closed
		}
		if f.Allow(port, et) {
			allowed++
		}
	}
	if allowed > 2 {
		t.Fatalf("expected at most 2 allowed events, got %d", allowed)
	}
}

func TestAllow_WindowResetAllowsAgain(t *testing.T) {
	now := time.Now()
	f := New(100*time.Millisecond, 1)
	f.now = func() time.Time { return now }

	port := uint16(22)
	f.Allow(port, Opened)
	f.Allow(port, Closed) // flap — count=2, suppressed

	// advance past window
	f.now = func() time.Time { return now.Add(200 * time.Millisecond) }
	if !f.Allow(port, Opened) {
		t.Fatal("expected event after window reset to be allowed")
	}
}

func TestAllow_DifferentPortsAreIndependent(t *testing.T) {
	f := New(time.Minute, 1)
	f.Allow(80, Opened)
	f.Allow(80, Closed) // flap on 80

	if !f.Allow(443, Opened) {
		t.Fatal("port 443 should be independent of port 80 flap")
	}
}

func TestAllow_SameEventTypeNotCounted(t *testing.T) {
	f := New(time.Minute, 1)
	port := uint16(3000)
	f.Allow(port, Opened)
	// repeated Opened events should not increment flap count
	if !f.Allow(port, Opened) {
		t.Fatal("repeated same-type event should not be suppressed")
	}
}
