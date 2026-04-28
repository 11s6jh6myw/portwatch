package portmeta

import (
	"github.com/user/portwatch/internal/scanner"
)

// NewProminenceAnnotator returns an annotator that attaches a prominence level
// to each port's metadata under the key "prominence".
func NewProminenceAnnotator() func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			if p.Meta == nil {
				p.Meta = make(map[string]string)
			}
			level := ProminenceFor(p)
			p.Meta["prominence"] = level.String()
			out[i] = p
		}
		return out
	}
}

// FilterByMinProminence returns only ports whose recorded prominence is at or
// above the given minimum level. Ports with no metadata pass through unchanged.
func FilterByMinProminence(ports []scanner.PortInfo, min ProminenceLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if p.Meta == nil {
			out = append(out, p)
			continue
		}
		raw, ok := p.Meta["prominence"]
		if !ok {
			out = append(out, p)
			continue
		}
		if parseProminence(raw) >= min {
			out = append(out, p)
		}
	}
	return out
}

func parseProminence(s string) ProminenceLevel {
	switch s {
	case "low":
		return ProminenceLow
	case "medium":
		return ProminenceMedium
	case "high":
		return ProminenceHigh
	case "critical":
		return ProminenceCritical
	default:
		return ProminenceNone
	}
}
