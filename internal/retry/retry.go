package retry

import (
	"context"
	"time"
)

// Policy defines how retries are attempted.
type Policy struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// Do executes fn up to MaxAttempts times, backing off between attempts.
// It stops early if ctx is cancelled or fn returns nil.
func Do(ctx context.Context, p Policy, fn func() error) error {
	delay := p.Delay
	var lastErr error
	for attempt := 0; attempt < p.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if attempt < p.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * p.Multiplier)
		}
	}
	return lastErr
}
