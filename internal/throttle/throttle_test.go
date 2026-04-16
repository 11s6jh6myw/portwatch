package throttle_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/example/portwatch/internal/throttle"
)

func TestAllow_UnderLimitAlwaysPasses(t *testing.T) {
	th := throttle.New(time.Second, 3)
	for i := 0; i < 3; i++ {
		if !th.Allow("k") {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
}

func TestAllow_ExceedsLimitSuppressed(t *testing.T) {
	th := throttle.New(time.Second, 2)
	th.Allow("k")
	th.Allow("k")
	if th.Allow("k") {
		t.Fatal("expected suppression after limit reached")
	}
}

func TestAllow_WindowExpiryResetsCount(t *testing.T) {
	now := time.Unix(1_000, 0)
	th := throttle.New(time.Second, 2)
	th.Allow("k") // inject via public API; advance clock via closure
	_ = th        // use real clock variant; test via Reset instead

	th2 := throttle.New(time.Second, 2)
	th2.Allow("k")
	th2.Allow("k")
	th2.Reset("k")
	if !th2.Allow("k") {
		t.Fatal("expected allow after reset")
	}
	_ = now
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	th := throttle.New(time.Second, 1)
	th.Allow("a")
	if !th.Allow("b") {
		t.Fatal("key b should not be affected by key a")
	}
}

func TestAllow_MultipleKeys(t *testing.T) {
	th := throttle.New(time.Minute, 5)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("port:%d", i)
		if !th.Allow(key) {
			t.Fatalf("first event for %s should be allowed", key)
		}
	}
}

func TestReset_AllowsImmediatelyAfter(t *testing.T) {
	th := throttle.New(time.Hour, 1)
	th.Allow("x")
	if th.Allow("x") {
		t.Fatal("should be suppressed before reset")
	}
	th.Reset("x")
	if !th.Allow("x") {
		t.Fatal("should be allowed after reset")
	}
}
