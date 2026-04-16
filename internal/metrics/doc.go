// Package metrics provides runtime observability for the portwatch daemon.
//
// A Collector accumulates counters and timing information as the daemon
// scans ports and emits alerts. Metrics are exposed over a lightweight
// HTTP server with two endpoints:
//
//	/metrics  – JSON snapshot of current counters and timing.
//	/healthz  – Simple liveness probe that always returns 200 OK.
//
// Typical usage:
//
//	c := metrics.New()
//	go metrics.NewServer(":9090", c).ListenAndServe()
//
//	// inside scan loop:
//	start := time.Now()
//	ports := scanner.Scan()
//	c.RecordScan(time.Since(start), len(ports))
package metrics
