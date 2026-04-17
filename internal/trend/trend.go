// Package trend analyses port scan history to detect recurring patterns
// and flag ports that open/close frequently within a time window.
package trend

import (
	"sync"
	"time"
)

// Event represents a single open/close transition for a port.
type Event struct {
	Port      uint16
	Kind      string // "opened" | "closed"
	Timestamp time.Time
}

// PortTrend summarises activity for a single port.
type PortTrend struct {
	Port        uint16
	OpenCount   int
	CloseCount  int
	FirstSeen   time.Time
	LastSeen    time.Time
	Flapping    bool
}

// Analyzer tracks port events and computes trends.
type Analyzer struct {
	mu       sync.Mutex
	events   []Event
	window   time.Duration
	flapMin  int
}

// New returns an Analyzer that considers events within window and marks a
// port as flapping when total transitions >= flapMin.
func New(window time.Duration, flapMin int) *Analyzer {
	return &Analyzer{window: window, flapMin: flapMin}
}

// Record adds an event to the analyzer.
func (a *Analyzer) Record(e Event) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = append(a.events, e)
}

// Analyze returns a trend summary for every port seen within the window.
func (a *Analyzer) Analyze(now time.Time) []PortTrend {
	a.mu.Lock()
	defer a.mu.Unlock()

	cutoff := now.Add(-a.window)
	type counts struct {
		opened, closed int
		first, last    time.Time
	}
	m := map[uint16]*counts{}

	for _, e := range a.events {
		if e.Timestamp.Before(cutoff) {
			continue
		}
		c, ok := m[e.Port]
		if !ok {
			c = &counts{first: e.Timestamp, last: e.Timestamp}
			m[e.Port] = c
		}
		if e.Timestamp.Before(c.first) {
			c.first = e.Timestamp
		}
		if e.Timestamp.After(c.last) {
			c.last = e.Timestamp
		}
		if e.Kind == "opened" {
			c.opened++
		} else {
			c.closed++
		}
	}

	trends := make([]PortTrend, 0, len(m))
	for port, c := range m {
		total := c.opened + c.closed
		trends = append(trends, PortTrend{
			Port:       port,
			OpenCount:  c.opened,
			CloseCount: c.closed,
			FirstSeen:  c.first,
			LastSeen:   c.last,
			Flapping:   total >= a.flapMin,
		})
	}
	return trends
}
