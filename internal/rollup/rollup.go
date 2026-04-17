// Package rollup batches multiple port events within a time window
// and emits a single aggregated summary, reducing alert noise during
// large network changes.
package rollup

import (
	"sync"
	"time"

	"github.com/joshbeard/portwatch/internal/scanner"
)

// Event represents a single port change event.
type Event struct {
	Port   scanner.PortInfo
	Opened bool // true = opened, false = closed
}

// Summary is the aggregated result emitted after a window closes.
type Summary struct {
	Opened []scanner.PortInfo
	Closed []scanner.PortInfo
	At     time.Time
}

// Roller collects events within a window then flushes them.
type Roller struct {
	mu      sync.Mutex
	window  time.Duration
	events  []Event
	timer   *time.Timer
	onFlush func(Summary)
}

// New creates a Roller that waits window duration after the first event
// before calling onFlush with the aggregated Summary.
func New(window time.Duration, onFlush func(Summary)) *Roller {
	return &Roller{window: window, onFlush: onFlush}
}

// Add adds an event to the current batch, starting the flush timer if needed.
func (r *Roller) Add(e Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, e)
	if r.timer == nil {
		r.timer = time.AfterFunc(r.window, r.flush)
	}
}

func (r *Roller) flush() {
	r.mu.Lock()
	events := r.events
	r.events = nil
	r.timer = nil
	r.mu.Unlock()

	s := Summary{At: time.Now()}
	for _, e := range events {
		if e.Opened {
			s.Opened = append(s.Opened, e.Port)
		} else {
			s.Closed = append(s.Closed, e.Port)
		}
	}
	r.onFlush(s)
}

// Flush forces an immediate flush regardless of the window timer.
func (r *Roller) Flush() {
	r.mu.Lock()
	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
	r.mu.Unlock()
	r.flush()
}
