// Package portmeta — stability.go
//
// StabilityLevel classifies how consistently a port has remained open over
// time. It combines the port's age (first-seen timestamp) with the number of
// open/close transitions recorded in history.
//
// Levels (lowest → highest):
//
//	unknown   – no first-seen data available
//	unstable  – 10+ state changes
//	variable  – 4–9 state changes
//	stable    – fewer than 4 changes
//	locked    – open for 30+ days with at most one change
//
// The annotator reads "stability.first_seen" (RFC3339) and
// "stability.changes" (integer string) from port metadata and writes
// "stability" back as a human-readable level string.
package portmeta
