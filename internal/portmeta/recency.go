package portmeta

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// RecencyLevel describes how recently a port was last observed open.
type RecencyLevel int

const (
	RecencyNone    RecencyLevel = iota // never seen or no data
	RecencyStale                       // last seen > 7 days ago
	RecencyAged                        // last seen 1–7 days ago
	RecencyRecent                      // last seen within 24 hours
	RecencyLive                        // last seen within 5 minutes
)

func (r RecencyLevel) String() string {
	switch r {
	case RecencyStale:
		return "stale"
	case RecencyAged:
		return "aged"
	case RecencyRecent:
		return "recent"
	case RecencyLive:
		return "live"
	default:
		return "none"
	}
}

// RecencyFor returns a RecencyLevel based on the port's last-seen timestamp
// stored in its metadata. If no timestamp is present, RecencyNone is returned.
func RecencyFor(p scanner.PortInfo) RecencyLevel {
	if p.Meta == nil {
		return RecencyNone
	}
	raw, ok := p.Meta["last_seen"]
	if !ok {
		return RecencyNone
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return RecencyNone
	}
	age := time.Since(t)
	switch {
	case age <= 5*time.Minute:
		return RecencyLive
	case age <= 24*time.Hour:
		return RecencyRecent
	case age <= 7*24*time.Hour:
		return RecencyAged
	default:
		return RecencyStale
	}
}

// IsLive returns true when the port was observed within the last 5 minutes.
func IsLive(p scanner.PortInfo) bool {
	return RecencyFor(p) == RecencyLive
}
