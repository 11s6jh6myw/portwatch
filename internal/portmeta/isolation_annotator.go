package portmeta

import "github.com/iamcalledrob/portwatch/internal/scanner"

const isolationKey = "isolation"

// NewIsolationAnnotator returns an annotator that sets the "isolation" metadata
// key on each port based on its peer set.
func NewIsolationAnnotator(peers []scanner.PortInfo) func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			level := IsolationFor(p, peers)
			if p.Meta == nil {
				p.Meta = make(map[string]string)
			}
			p.Meta[isolationKey] = level.String()
			out[i] = p
		}
		return out
	}
}

// FilterByMinIsolation returns only ports whose recorded isolation level is at
// least min. Ports with no isolation metadata are passed through unchanged.
func FilterByMinIsolation(ports []scanner.PortInfo, min IsolationLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if p.Meta == nil {
			out = append(out, p)
			continue
		}
		level := parseIsolation(p.Meta[isolationKey])
		if level >= min {
			out = append(out, p)
		}
	}
	return out
}

func parseIsolation(s string) IsolationLevel {
	switch s {
	case "low":
		return IsolationLow
	case "medium":
		return IsolationMedium
	case "high":
		return IsolationHigh
	default:
		return IsolationNone
	}
}
