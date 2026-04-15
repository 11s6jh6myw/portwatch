// Package ratelimit provides a thread-safe rate limiter for port change alerts.
//
// It prevents alert fatigue by suppressing repeated notifications for the same
// port event (e.g. "tcp/8080 opened") within a configurable cooldown window.
//
// Each unique combination of protocol, port number, and event type is tracked
// independently. An alert is allowed through only if no alert for the same
// combination has been issued within the cooldown duration.
//
// Usage:
//
//	limiter := ratelimit.New(30 * time.Second)
//
//	if limiter.Allow("tcp", 8080, "opened") {
//		// send alert
//	}
//
// Call Expire periodically (e.g. in a background goroutine or after each scan
// cycle) to reclaim memory from entries whose cooldown window has elapsed:
//
//	limiter.Expire()
//
// The zero value of Limiter is not usable; always construct with New.
package ratelimit
