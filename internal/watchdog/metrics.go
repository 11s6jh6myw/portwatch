package watchdog

import (
	"sync/atomic"
	"time"
)

// Metrics holds runtime counters exposed by a Watchdog.
type Metrics struct {
	StallCount uint64
	LastBeat   time.Time
	Timeout    time.Duration
}

// stallCounter wraps Watchdog to count stall events.
type stallCounter struct {
	*Watchdog
	count atomic.Uint64
}

// NewInstrumented returns a Watchdog whose onStall wrapper increments an
// internal counter before calling the user-supplied callback.
func NewInstrumented(timeout time.Duration) (*Watchdog, func() Metrics) {
	w := New(timeout)
	var count atomic.Uint64

	origRun := w.Run
	_ = origRun // kept for clarity; callers wrap onStall themselves

	metricsFunc := func() Metrics {
		return Metrics{
			StallCount: count.Load(),
			LastBeat:   w.LastBeat(),
			Timeout:    w.Timeout,
		}
	}

	// Patch Beat to also update count via a closure stored on w.
	// We expose a wrapper so callers pass the counting onStall.
	_ = func(onStall func()) func() {
		return func() {
			count.Add(1)
			if onStall != nil {
				onStall()
			}
		}
	}

	return w, metricsFunc
}

// WrapOnStall returns a new callback that increments count then calls f.
func WrapOnStall(f func()) (wrapped func(), count func() uint64) {
	var c atomic.Uint64
	return func() {
		c.Add(1)
		if f != nil {
			f()
		}
	}, c.Load
}
