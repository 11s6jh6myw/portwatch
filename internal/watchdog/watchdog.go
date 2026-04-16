// Package watchdog provides a self-monitoring loop that restarts the
// scan cycle if it stalls beyond a configurable deadline.
package watchdog

import (
	"context"
	"log"
	"sync"
	"time"
)

// Watchdog monitors a heartbeat channel and cancels the supplied context
// if no beat is received within Timeout.
type Watchdog struct {
	Timeout   time.Duration
	heartbeat chan struct{}
	mu        sync.Mutex
	last      time.Time
}

// New creates a Watchdog with the given timeout.
func New(timeout time.Duration) *Watchdog {
	return &Watchdog{
		Timeout:   timeout,
		heartbeat: make(chan struct{}, 1),
		last:      time.Now(),
	}
}

// Beat signals that the monitored process is alive.
func (w *Watchdog) Beat() {
	w.mu.Lock()
	w.last = time.Now()
	w.mu.Unlock()
	select {
	case w.heartbeat <- struct{}{}:
	default:
	}
}

// Run starts the watchdog loop. It calls onStall whenever a timeout is
// detected, then waits for the next beat before watching again.
// It returns when ctx is cancelled.
func (w *Watchdog) Run(ctx context.Context, onStall func()) {
	ticker := time.NewTicker(w.Timeout / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.mu.Lock()
			since := time.Since(w.last)
			w.mu.Unlock()
			if since > w.Timeout {
				log.Printf("watchdog: stall detected (last beat %s ago)", since.Round(time.Millisecond))
				onStall()
			}
		}
	}
}

// LastBeat returns the time of the most recent heartbeat.
func (w *Watchdog) LastBeat() time.Time {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.last
}
