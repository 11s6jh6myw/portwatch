// Package alert provides alerting primitives for portwatch.
//
// A Notifier emits human-readable, timestamped log lines whenever a port
// change event is detected by the monitor. Each line follows the format:
//
//	[<RFC3339 timestamp>] <LEVEL> <proto>/<port> — <message>
//
// Example output:
//
//	[2024-05-01T12:00:00Z] ALERT tcp/8080 — port opened
//	[2024-05-01T12:01:00Z] WARN  tcp/22   — port closed
//
// Usage:
//
//	n := alert.New(os.Stderr)
//	n.NotifyOpened(portInfo)
//	n.NotifyClosed(portInfo)
package alert
