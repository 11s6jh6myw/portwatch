// Package ratelimit provides a thread-safe rate limiter for port change alerts.
//
// It prevents alert fatigue by suppressing repeated notifications for the same
// port event (e.g. "tcp/8080 opened") within a configurable cooldown window.
//
// Usage:
//
//	limiter := ratelimit.New(30 * time.Second)
//
//	if limiter.Allow("tcp", 8080, "opened") {
//		// send alert
//	}
//
// Call Expire periodically to reclaim memory from stale entries.
package ratelimit
