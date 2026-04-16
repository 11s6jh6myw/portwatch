// Package suppress provides a flap-detection filter that suppresses
// alerts for ports that open and close repeatedly within a short window.
package suppress

import (
	"sync"
	"time"
)

// EventType distinguishes an open from a close event.
type EventType int

const (
	Opened EventType = iota
	Closed
)

type entry struct {
	count     int
	firstSeen time.Time
	lastType  EventType
}

// Filter tracks port state-change frequency and suppresses alerts when a
// port flaps more than Threshold times within Window.
type Filter struct {
	mu        sync.Mutex
	Window    time.Duration
	Threshold int
	entries   map[string]*entry
	now       func() time.Time
}

// New returns a Filter with the given flap window and threshold.
func New(window time.Duration, threshold int) *Filter {
	return &Filter{
		Window:    window,
		Threshold: threshold,
		entries:   make(map[string]*entry),
		now:       time.Now,
	}
}

// Allow returns true if the event should be forwarded (not suppressed).
// It records every call and suppresses once the port has flapped more
// than Threshold times inside the current Window.
func (f *Filter) Allow(port uint16, et EventType) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	key := portKey(port)
	now := f.now()

	e, ok := f.entries[key]
	if !ok || now.Sub(e.firstSeen) > f.Window {
		f.entries[key] = &entry{count: 1, firstSeen: now, lastType: et}
		return true
	}

	if et != e.lastType {
		e.count++
		e.lastType = et
	}

	return e.count <= f.Threshold
}

func portKey(port uint16) string {
	buf := [5]byte{}
	buf[0] = byte(port >> 8)
	buf[1] = byte(port)
	return string(buf[:2])
}
