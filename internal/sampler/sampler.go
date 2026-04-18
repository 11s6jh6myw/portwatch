// Package sampler provides periodic port-scan sampling with jitter to avoid
// thundering-herd effects when multiple portwatch instances run together.
package sampler

import (
	"context"
	"math/rand"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Result holds a single sampled scan outcome.
type Result struct {
	Ports []scanner.PortInfo
	At    time.Time
	Err   error
}

// Sampler runs a scanner on a fixed interval with optional jitter.
type Sampler struct {
	scanner  *scanner.TCPScanner
	interval time.Duration
	jitter   time.Duration
	results  chan Result
}

// New creates a Sampler. jitter adds a random delay in [0, jitter) before
// each scan to spread load.
func New(s *scanner.TCPScanner, interval, jitter time.Duration) *Sampler {
	return &Sampler{
		scanner:  s,
		interval: interval,
		jitter:   jitter,
		results:  make(chan Result, 4),
	}
}

// Results returns the read-only channel of scan results.
func (s *Sampler) Results() <-chan Result { return s.results }

// Run starts the sampling loop and blocks until ctx is cancelled.
func (s *Sampler) Run(ctx context.Context) {
	defer close(s.results)
	for {
		s.sleepWithJitter(ctx)
		if ctx.Err() != nil {
			return
		}
		ports, err := s.scanner.Scan(ctx)
		select {
		case s.results <- Result{Ports: ports, At: time.Now(), Err: err}:
		case <-ctx.Done():
			return
		}
		select {
		case <-time.After(s.interval):
		case <-ctx.Done():
			return
		}
	}
}

func (s *Sampler) sleepWithJitter(ctx context.Context) {
	if s.jitter <= 0 {
		return
	}
	d := time.Duration(rand.Int63n(int64(s.jitter)))
	select {
	case <-time.After(d):
	case <-ctx.Done():
	}
}
