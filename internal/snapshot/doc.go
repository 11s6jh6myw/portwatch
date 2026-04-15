// Package snapshot defines the Snapshot type used to capture a point-in-time
// view of open ports detected on the host.
//
// A Snapshot records the hostname, UTC timestamp, and list of open ports.
// It provides helpers for equality comparison and fast port lookup, making it
// straightforward to detect changes between successive scans.
//
// Typical usage:
//
//	ports, _ := scanner.Scan(cfg)
//	snap := snapshot.New(hostname, ports)
//	if !snap.Equal(previous) {
//		// handle change
//	}
package snapshot
