// Package cooldown provides a concurrency-safe per-key cooldown tracker.
//
// It is useful for suppressing repeated alerts or actions for the same
// event within a configurable quiet period. Unlike ratelimit which counts
// events per window, cooldown simply enforces a minimum gap between
// successive allowed actions for each key.
//
// Example:
//
//	tr := cooldown.New(30 * time.Second)
//	if tr.Allow(portKey) {
//		sendAlert(portKey)
//	}
package cooldown
