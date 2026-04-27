package portmeta

import (
	"github.com/user/portwatch/internal/scanner"
)

// NewRecencyAnnotator returns an annotator that writes the recency level and
// the is_live flag into each port's metadata map.
func NewRecencyAnnotator() func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		for i := range ports {
			if ports[i].Meta == nil {
				ports[i].Meta = make(map[string]string)
			}
			level := RecencyFor(ports[i])
			ports[i].Meta["recency"] = level.String()
			if level == RecencyLive {
				ports[i].Meta["is_live"] = "true"
			} else {
				ports[i].Meta["is_live"] = "false"
			}
		}
		return ports
	}
}

// FilterByMinRecency returns only ports whose recency is at least min.
func FilterByMinRecency(ports []scanner.PortInfo, min RecencyLevel) []scanner.PortInfo {
	out := make([]scanner.PortInfo, 0, len(ports))
	for _, p := range ports {
		if RecencyFor(p) >= min {
			out = append(out, p)
		}
	}
	return out
}

// parseRecency converts a string to a RecencyLevel, defaulting to RecencyNone.
func parseRecency(s string) RecencyLevel {
	switch s {
	case "live":
		return RecencyLive
	case "recent":
		return RecencyRecent
	case "aged":
		return RecencyAged
	case "stale":
		return RecencyStale
	default:
		return RecencyNone
	}
}
