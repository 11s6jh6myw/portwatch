// Package eventlog provides an append-only, newline-delimited JSON log of
// port change events (opened / closed) detected by portwatch.
//
// Events are written atomically under a mutex so concurrent monitor goroutines
// can safely call Append without data races. ReadAll and Query replay the full
// log from disk, making the store suitable for audit trails and post-hoc
// analysis rather than high-frequency hot paths.
package eventlog
