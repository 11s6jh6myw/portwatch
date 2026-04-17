// Package cooldown provides a per-key cooldown tracker that prevents
// repeated actions within a configurable quiet period.
package cooldown

import (
	"sync"
	"time"
)

// Tracker tracks the last action time per key and reports whether
// enough time has elapsed to allow the next action.
type Tracker struct {
	mu      sync.Mutex
	period  time.Duration
	last    map[string]time.Time
	nowFunc func() time.Time
}

// New returns a Tracker with the given cooldown period.
func New(period time.Duration) *Tracker {
	return &Tracker{
		period:  period,
		last:    make(map[string]time.Time),
		nowFunc: time.Now,
	}
}

// Allow returns true if the cooldown period has elapsed since the last
// call to Allow for the given key, and records the current time.
// The first call for any key always returns true.
func (t *Tracker) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.nowFunc()
	if last, ok := t.last[key]; ok && now.Sub(last) < t.period {
		return false
	}
	t.last[key] = now
	return true
}

// Reset removes the cooldown record for the given key so the next
// call to Allow will always return true.
func (t *Tracker) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// Len returns the number of keys currently tracked.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.last)
}
