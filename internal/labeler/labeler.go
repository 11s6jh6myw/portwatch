// Package labeler assigns human-readable labels to ports based on
// well-known service mappings and user-defined overrides.
package labeler

import "fmt"

// well-known maps port numbers to service names.
var wellKnown = map[uint16]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Labeler resolves a label for a given port number.
type Labeler struct {
	overrides map[uint16]string
}

// New returns a Labeler with optional user-defined overrides.
// Override entries take precedence over well-known mappings.
func New(overrides map[uint16]string) *Labeler {
	if overrides == nil {
		overrides = make(map[uint16]string)
	}
	return &Labeler{overrides: overrides}
}

// Label returns the service label for port. If no mapping exists the
// label is the numeric port formatted as a string.
func (l *Labeler) Label(port uint16) string {
	if v, ok := l.overrides[port]; ok {
		return v
	}
	if v, ok := wellKnown[port]; ok {
		return v
	}
	return fmt.Sprintf("%d", port)
}

// Known reports whether port has a recognised label (override or
// well-known).
func (l *Labeler) Known(port uint16) bool {
	_, inOverride := l.overrides[port]
	_, inWellKnown := wellKnown[port]
	return inOverride || inWellKnown
}
