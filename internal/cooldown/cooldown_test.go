package cooldown

import (
	"testing"
	"time"
)

func TestAllow_FirstCallAlwaysTrue(t *testing.T) {
	tr := New(time.Second)
	if !tr.Allow("k1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_WithinCooldownSuppressed(t *testing.T) {
	tr := New(time.Second)
	tr.Allow("k1")
	if tr.Allow("k1") {
		t.Fatal("expected second call within cooldown to be suppressed")
	}
}

func TestAllow_AfterCooldownAllowed(t *testing.T) {
	now := time.Now()
	tr := New(time.Second)
	tr.nowFunc = func() time.Time { return now }
	tr.Allow("k1")
	tr.nowFunc = func() time.Time { return now.Add(2 * time.Second) }
	if !tr.Allow("k1") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	tr := New(time.Second)
	tr.Allow("k1")
	if !tr.Allow("k2") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestReset_AllowsImmediately(t *testing.T) {
	tr := New(time.Second)
	tr.Allow("k1")
	tr.Reset("k1")
	if !tr.Allow("k1") {
		t.Fatal("expected allow after reset")
	}
}

func TestLen_TracksKeys(t *testing.T) {
	tr := New(time.Second)
	if tr.Len() != 0 {
		t.Fatalf("expected 0, got %d", tr.Len())
	}
	tr.Allow("a")
	tr.Allow("b")
	if tr.Len() != 2 {
		t.Fatalf("expected 2, got %d", tr.Len())
	}
	tr.Reset("a")
	if tr.Len() != 1 {
		t.Fatalf("expected 1, got %d", tr.Len())
	}
}
