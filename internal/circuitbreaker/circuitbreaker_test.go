package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/circuitbreaker"
)

func TestAllow_ClosedByDefault(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllow_OpensAfterThreshold(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if err := b.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestAllow_SuccessResetsFailures(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	b.RecordSuccess()
	b.RecordFailure()
	if b.CurrentState() != circuitbreaker.StateClosed {
		t.Fatal("expected circuit to remain closed after success reset")
	}
}

func TestAllow_ReopensDuringCooldown(t *testing.T) {
	b := circuitbreaker.New(2, 5*time.Second)
	b.RecordFailure()
	b.RecordFailure()
	if err := b.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen during cooldown, got %v", err)
	}
}

func TestAllow_RecoverAfterCooldown(t *testing.T) {
	now := time.Now()
	b := circuitbreaker.New(2, 10*time.Millisecond)
	// inject controllable clock via unexported field workaround: use real sleep
	b.RecordFailure()
	b.RecordFailure()

	time.Sleep(20 * time.Millisecond)
	_ = now // suppress unused

	if err := b.Allow(); err != nil {
		t.Fatalf("expected circuit to recover after cooldown, got %v", err)
	}
	if b.CurrentState() != circuitbreaker.StateClosed {
		t.Fatal("expected StateClosed after recovery")
	}
}

func TestCurrentState_ReflectsTransitions(t *testing.T) {
	b := circuitbreaker.New(1, time.Second)
	if b.CurrentState() != circuitbreaker.StateClosed {
		t.Fatal("initial state should be closed")
	}
	b.RecordFailure()
	if b.CurrentState() != circuitbreaker.StateOpen {
		t.Fatal("state should be open after threshold")
	}
	b.RecordSuccess()
	if b.CurrentState() != circuitbreaker.StateClosed {
		t.Fatal("state should be closed after success")
	}
}
