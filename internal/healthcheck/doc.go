// Package healthcheck implements a lightweight TCP probe loop used to
// verify that the portwatch daemon itself can reach the local network
// stack. When a probe fails the registered callback is invoked so the
// caller can emit an alert, increment a metric, or trigger a watchdog
// beat. Typical usage:
//
//	checker := healthcheck.New(
//		"127.0.0.1:80",
//		30*time.Second,
//		2*time.Second,
//		func(s healthcheck.Status) { log.Printf("health probe failed: %v", s.Err) },
//	)
//	go checker.Run(ctx)
package healthcheck
