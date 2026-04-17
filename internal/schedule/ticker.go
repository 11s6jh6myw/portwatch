package schedule

import (
	"context"
	"time"
)

// Ticker fires at the adaptive interval managed by a Scheduler.
// Each tick delivers the current time on C; after reading a tick the caller
// must call either RecordActivity or RecordQuiet so the interval adapts.
type Ticker struct {
	C         <-chan time.Time
	scheduler *Scheduler
	stop      chan struct{}
}

// NewTicker starts an adaptive ticker driven by s.
func NewTicker(ctx context.Context, s *Scheduler) *Ticker {
	ch := make(chan time.Time, 1)
	t := &Ticker{C: ch, scheduler: s, stop: make(chan struct{})}
	go t.run(ctx, ch)
	return t
}

func (t *Ticker) run(ctx context.Context, ch chan<- time.Time) {
	for {
		interval := t.scheduler.Current()
		select {
		case <-time.After(interval):
			select {
			case ch <- time.Now():
			default:
			}
		case <-ctx.Done():
			return
		case <-t.stop:
			return
		}
	}
}

// Stop terminates the ticker goroutine.
func (t *Ticker) Stop() {
	close(t.stop)
}
