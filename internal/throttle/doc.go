// Package throttle provides a sliding-window rate limiter for alert events.
//
// A Throttle tracks how many times each key has triggered within a configurable
// time window and suppresses further events once the per-key limit is reached.
// The window slides forward continuously; old timestamps are evicted lazily on
// each call to Allow.
//
// Typical usage:
//
//	th := throttle.New(time.Minute, 5)
//	if th.Allow(eventKey) {
//	    // dispatch notification
//	}
package throttle
