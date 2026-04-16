// Package watchdog implements a heartbeat-based stall detector for the
// portwatch scan loop.
//
// The Watchdog expects periodic calls to Beat; if none arrive within the
// configured Timeout the supplied onStall callback is invoked. This lets
// the main daemon restart or log a warning when the scanner hangs (e.g.
// due to a network issue or a blocked goroutine).
//
// Typical usage:
//
//	wd := watchdog.New(10 * time.Second)
//	go wd.Run(ctx, func() {
//	    log.Println("scan loop appears stalled — restarting")
//	    restartScan()
//	})
//	// inside the scan loop:
//	wd.Beat()
package watchdog
