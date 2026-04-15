// Package reporter provides formatted output of port change reports.
//
// It supports two output formats:
//
//   - text: human-readable lines prefixed with + (opened) or - (closed)
//   - json: machine-readable JSON objects, one per line
//
// Usage:
//
//	r := reporter.New(os.Stdout, reporter.FormatText)
//	err := r.Write(reporter.Report{
//	    Timestamp: time.Now(),
//	    Opened:    openedPorts,
//	    Closed:    closedPorts,
//	})
package reporter
