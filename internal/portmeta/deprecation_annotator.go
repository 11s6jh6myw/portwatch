package portmeta

import "github.com/user/portwatch/internal/scanner"

const (
	metaDeprecationLevel = "deprecation.level"
	metaDeprecationScore = "deprecation.score"
)

// NewDeprecationAnnotator returns an annotator that adds deprecation metadata
// to each PortInfo using the port's known legacy status.
func NewDeprecationAnnotator() func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			level := DeprecationFor(p)
			if p.Meta == nil {
				p.Meta = make(map[string]string)
			}
			p.Meta[metaDeprecationLevel] = level.String()
			p.Meta[metaDeprecationScore] = itoa(int(level))
			out[i] = p
		}
		return out
	}
}

// FilterByMaxDeprecation returns only ports whose deprecation level is at or
// below the given maximum. Ports without deprecation metadata are passed through.
func FilterByMaxDeprecation(ports []scanner.PortInfo, max DeprecationLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if p.Meta == nil {
			out = append(out, p)
			continue
		}
		v, ok := p.Meta[metaDeprecationLevel]
		if !ok {
			out = append(out, p)
			continue
		}
		if parseDeprecation(v) <= max {
			out = append(out, p)
		}
	}
	return out
}

func parseDeprecation(s string) DeprecationLevel {
	switch s {
	case "minor":
		return DeprecationMinor
	case "moderate":
		return DeprecationModerate
	case "high":
		return DeprecationHigh
	case "obsolete":
		return DeprecationObsolete
	default:
		return DeprecationNone
	}
}
