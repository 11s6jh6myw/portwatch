package portmeta

import "github.com/iamcalledrob/portwatch/internal/scanner"

// NewAffinityAnnotator returns an annotator that adds "affinity" and
// "affinity.family" metadata keys to each PortInfo.
func NewAffinityAnnotator() func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			level, family := AffinityFor(p)
			if p.Meta == nil {
				p.Meta = make(map[string]string)
			}
			p.Meta["affinity"] = level.String()
			if family != "" {
				p.Meta["affinity.family"] = family
			}
			out[i] = p
		}
		return out
	}
}

// FilterByMinAffinity returns only ports whose affinity level is >= min.
// Ports with no metadata are passed through unchanged.
func FilterByMinAffinity(ports []scanner.PortInfo, min AffinityLevel) []scanner.PortInfo {
	out := ports[:0:0]
	for _, p := range ports {
		if p.Meta == nil {
			out = append(out, p)
			continue
		}
		raw, ok := p.Meta["affinity"]
		if !ok {
			out = append(out, p)
			continue
		}
		level := parseAffinity(raw)
		if level >= min {
			out = append(out, p)
		}
	}
	return out
}

func parseAffinity(s string) AffinityLevel {
	switch s {
	case "weak":
		return AffinityWeak
	case "medium":
		return AffinityMedium
	case "strong":
		return AffinityStrong
	default:
		return AffinityNone
	}
}
