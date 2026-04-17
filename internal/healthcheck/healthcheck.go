// Package healthcheck provides a periodic self-test that verifies the
// scanner can reach localhost, emitting a synthetic event on failure.
package healthcheck

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Status represents the result of a single health probe.
type Status struct {
	Healthy   bool
	CheckedAt time.Time
	Err       error
}

// Checker runs periodic TCP probes against a known-open port.
type Checker struct {
	addr     string
	interval time.Duration
	timeout  time.Duration
	onFail   func(Status)
}

// New creates a Checker that dials addr every interval.
// onFail is called whenever a probe fails.
func New(addr string, interval, timeout time.Duration, onFail func(Status)) *Checker {
	return &Checker{
		addr:     addr,
		interval: interval,
		timeout:  timeout,
		onFail:   onFail,
	}
}

// Run starts the probe loop and blocks until ctx is cancelled.
func (c *Checker) Run(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s := c.probe()
			if !			}
 Status {
	deadline := time.Now().Add(c.timeout)
	conn, err := net.DialTimeout("tcp", c.addr, time.Until(deadline))
	if err != nil {
		return Status{Healthy: false, CheckedAt: time.Now(), Err: fmt.Errorf("dial %s: %w", c.addr, err)}
	}
	_ = conn.Close()
	return Status{Healthy: true, CheckedAt: time.Now()}
}
