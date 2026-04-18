// Package sampler provides periodic port-scan sampling with configurable
// intervals and jitter.
//
// # Overview
//
// A Sampler wraps a scanner.TCPScanner and emits Result values on a channel
// at a fixed cadence. A random jitter delay before each scan helps distribute
// load when multiple portwatch instances share the same network segment.
//
// # Aggregator
//
// NewAggregator fans multiple Samplers into a single output channel and
// exposes a Merge helper to deduplicate port lists from different hosts.
package sampler
