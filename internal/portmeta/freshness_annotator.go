package portmeta

import (
	"time"

	"github.com/iamcathal/portwatch/internal/scanner"
)

const freshnessKey = "freshness"

// FreshnessAnnotator attaches a freshness label to each port based on its
// first-seen timestamp stored in port metadata.
type FreshnessAnnotator struct {
	// firstSeen maps port number to the time it was first observed.
	firstSeen map[int]time.Time
}

// NewFreshnessAnnotator creates an annotator seeded with first-seen times.
func NewFreshnessAnnotator(firstSeen map[int]time.Time) *FreshnessAnnotator {
	if firstSeen == nil {
		firstSeen = make(map[int]time.Time)
	}
	return &FreshnessAnnotator{firstSeen: firstSeen}
}

// Annotate sets the "freshness" metadata key on each port.
func (a *FreshnessAnnotator) Annotate(ports []scanner.PortInfo) []scanner.PortInfo {
	out := make([]scanner.PortInfo, len(ports))
	for i, p := range ports {
		fs := ClassifyFreshness(a.firstSeen[p.Port])
		if p.Meta == nil {
			p.Meta = make(map[string]string)
		}
		p.Meta[freshnessKey] = fs.String()
		out[i] = p
	}
	return out
}

// FilterByMinFreshness returns only ports whose freshness is >= min.
func FilterByMinFreshness(ports []scanner.PortInfo, min FreshnessLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if p.Meta == nil {
			out = append(out, p)
			continue
		}
		lvl := parseFreshness(p.Meta[freshnessKey])
		if lvl >= min {
			out = append(out, p)
		}
	}
	return out
}

func parseFreshness(s string) FreshnessLevel {
	switch s {
	case "new":
		return FreshnessNew
	case "recent":
		return FreshnessRecent
	case "mature":
		return FreshnessMature
	case "stale":
		return FreshnessStale
	default:
		return FreshnessUnknown
	}
}
