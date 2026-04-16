// Package circuitbreaker provides a simple circuit breaker for alert delivery.
// It trips after a configurable number of consecutive failures and resets
// after a cooldown period, preventing alert storms against unavailable sinks.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit is open and calls are rejected.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current circuit breaker state.
type State int

const (
	StateClosed State = iota
	StateOpen
)

// Breaker is a simple circuit breaker.
type Breaker struct {
	mu           sync.Mutex
	failures     int
	threshold    int
	cooldown     time.Duration
	openedAt     time.Time
	state        State
	now          func() time.Time
}

// New creates a Breaker that opens after threshold consecutive failures
// and attempts recovery after cooldown.
func New(threshold int, cooldown time.Duration) *Breaker {
	return &Breaker{
		threshold: threshold,
		cooldown:  cooldown,
		now:       time.Now,
	}
}

// Allow reports whether a call should proceed.
// It returns ErrOpen when the circuit is open and the cooldown has not elapsed.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == StateOpen {
		if b.now().Sub(b.openedAt) >= b.cooldown {
			// half-open: allow one probe
			b.state = StateClosed
			b.failures = 0
		} else {
			return ErrOpen
		}
	}
	return nil
}

// RecordSuccess resets the failure counter.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure counter and opens the circuit if the
// threshold is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.threshold {
		b.state = StateOpen
		b.openedAt = b.now()
	}
}

// CurrentState returns the current state of the breaker.
func (b *Breaker) CurrentState() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
