// Package throttle limits the rate of outbound alerts by capping the number
// of notifications dispatched within a sliding time window.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks per-key event counts within a rolling window.
type Throttle struct {
	mu      sync.Mutex
	window  time.Duration
	maxRate int
	buckets map[string][]time.Time
	now     func() time.Time
}

// New returns a Throttle that allows at most maxRate events per key per window.
func New(window time.Duration, maxRate int) *Throttle {
	return &Throttle{
		window:  window,
		maxRate: maxRate,
		buckets: make(map[string][]time.Time),
		now:     time.Now,
	}
}

// Allow returns true if the event identified by key may proceed, false if it
// should be suppressed because the rate limit has been reached.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	cutoff := now.Add(-t.window)

	times := t.buckets[key]
	filtered := times[:0]
	for _, ts := range times {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}

	if len(filtered) >= t.maxRate {
		t.buckets[key] = filtered
		return false
	}

	t.buckets[key] = append(filtered, now)
	return true
}

// Reset clears the event history for the given key.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.buckets, key)
}
