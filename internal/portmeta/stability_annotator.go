package portmeta

import (
	"strconv"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

const (
	stabilityKey    = "stability"
	stabilitySeenKey = "stability.first_seen"
	stabilityChangesKey = "stability.changes"
)

// NewStabilityAnnotator returns an annotator that classifies port stability
// based on first-seen time and change count stored in port metadata.
func NewStabilityAnnotator() func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		for i, p := range ports {
			firstSeen := parseMetaTime(p.Meta[stabilitySeenKey])
			changes := parseMetaInt(p.Meta[stabilityChangesKey])
			level := ClassifyStability(firstSeen, changes)
			if ports[i].Meta == nil {
				ports[i].Meta = make(map[string]string)
			}
			ports[i].Meta[stabilityKey] = level.String()
		}
		return ports
	}
}

// FilterByMinStability returns only ports at or above the given stability level.
func FilterByMinStability(ports []scanner.PortInfo, min StabilityLevel) []scanner.PortInfo {
	out := ports[:0:0]
	for _, p := range ports {
		level := parseStability(p.Meta[stabilityKey])
		if level >= min {
			out = append(out, p)
		}
	}
	return out
}

func parseMetaTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func parseMetaInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func parseStability(s string) StabilityLevel {
	switch s {
	case "unstable":
		return StabilityUnstable
	case "variable":
		return StabilityVariable
	case "stable":
		return StabilityStable
	case "locked":
		return StabilityLocked
	default:
		return StabilityUnknown
	}
}
