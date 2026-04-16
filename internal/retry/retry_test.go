package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/retry"
)

var errFail = errors.New("fail")

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.DefaultPolicy(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	calls := 0
	p := retry.Policy{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1.0}
	err := retry.Do(context.Background(), p, func() error {
		calls++
		if calls < 3 {
			return errFail
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ReturnsLastError(t *testing.T) {
	p := retry.Policy{MaxAttempts: 2, Delay: time.Millisecond, Multiplier: 1.0}
	err := retry.Do(context.Background(), p, func() error {
		return errFail
	})
	if !errors.Is(err, errFail) {
		t.Fatalf("expected errFail, got %v", err)
	}
}

func TestDo_StopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	p := retry.Policy{MaxAttempts: 5, Delay: time.Millisecond, Multiplier: 1.0}
	err := retry.Do(ctx, p, func() error {
		calls++
		return errFail
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls > 1 {
		t.Fatalf("expected at most 1 call, got %d", calls)
	}
}
