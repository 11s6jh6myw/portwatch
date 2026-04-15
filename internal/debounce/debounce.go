// Package debounce provides a mechanism to suppress repeated events
// within a configurable time window, reducing alert noise from transient
// port state changes.
package debounce

import (
	"sync"
	"time"
)

// EventKey uniquely identifies a debounced event by port and event type.
type EventKey struct {
	Port  int
	Event string // "opened" or "closed"
}

// Debouncer suppresses repeated events that occur within a quiet window.
type Debouncer struct {
	mu      sync.Mutex
	window  time.Duration
	pending map[EventKey]time.Time
	clock   func() time.Time
}

// New creates a Debouncer that suppresses events repeated within window.
func New(window time.Duration) *Debouncer {
	return &Debouncer{
		window:  window,
		pending: make(map[EventKey]time.Time),
		clock:   time.Now,
	}
}

// Allow returns true if the event should be forwarded, or false if it
// falls within the debounce window of a previously seen identical event.
func (d *Debouncer) Allow(port int, event string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.evict()

	key := EventKey{Port: port, Event: event}
	if _, seen := d.pending[key]; seen {
		return false
	}

	d.pending[key] = d.clock().Add(d.window)
	return true
}

// evict removes expired entries from the pending map.
// Must be called with d.mu held.
func (d *Debouncer) evict() {
	now := d.clock()
	for k, expiry := range d.pending {
		if now.After(expiry) {
			delete(d.pending, k)
		}
	}
}

// Pending returns the number of events currently held in the debounce window.
func (d *Debouncer) Pending() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.evict()
	return len(d.pending)
}
