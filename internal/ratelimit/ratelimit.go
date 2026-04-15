// Package ratelimit provides alert rate limiting to suppress duplicate
// notifications for the same port event within a configurable cooldown window.
package ratelimit

import (
	"sync"
	"time"
)

// key uniquely identifies a port event by protocol, port number, and event type.
type key struct {
	protocol string
	port     int
	event    string
}

// Limiter suppresses repeated alerts for the same port event within a
// cooldown duration.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	seen     map[key]time.Time
	now      func() time.Time
}

// New returns a Limiter with the given cooldown window.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		seen:     make(map[key]time.Time),
		now:      time.Now,
	}
}

// Allow reports whether an alert for the given protocol, port, and event
// type should be emitted. It returns false if the same event was already
// allowed within the cooldown window.
func (l *Limiter) Allow(protocol string, port int, event string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	k := key{protocol: protocol, port: port, event: event}
	now := l.now()

	if last, ok := l.seen[k]; ok {
		if now.Sub(last) < l.cooldown {
			return false
		}
	}

	l.seen[k] = now
	return true
}

// Flush removes all recorded events, resetting the limiter state.
func (l *Limiter) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.seen = make(map[key]time.Time)
}

// Expire removes stale entries older than the cooldown to bound memory usage.
func (l *Limiter) Expire() {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	for k, t := range l.seen {
		if now.Sub(t) >= l.cooldown {
			delete(l.seen, k)
		}
	}
}
