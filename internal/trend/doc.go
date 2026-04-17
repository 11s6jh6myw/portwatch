// Package trend provides port activity trend analysis for portwatch.
//
// It tracks open/close events over a configurable time window and
// identifies ports exhibiting flapping behaviour — frequent transitions
// that may indicate unstable services or scanning activity.
//
// Usage:
//
//	a := trend.New(30*time.Minute, 4)
//	a.Record(trend.Event{Port: 443, Kind: "opened", Timestamp: time.Now()})
//	results := a.Analyze(time.Now())
package trend
