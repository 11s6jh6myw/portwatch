package portmeta

import (
	"strconv"

	"github.com/netwatch/portwatch/internal/scanner"
)

const (
	metaKeyCriticality = "criticality"
	metaKeyCriticalityScore = "criticality_score"
)

// NewCriticalityAnnotator returns an annotator that adds criticality metadata
// to each PortInfo based on its port number.
func NewCriticalityAnnotator() func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			c := CriticalityFor(uint16(p.Port))
			if p.Meta == nil {
				p.Meta = make(map[string]string)
			}
			p.Meta[metaKeyCriticality] = c.String()
			p.Meta[metaKeyCriticalityScore] = strconv.Itoa(int(c))
			out[i] = p
		}
		return out
	}
}

// FilterByMinCriticality returns only ports whose criticality meets the minimum.
func FilterByMinCriticality(ports []scanner.PortInfo, min CriticalityLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if CriticalityFor(uint16(p.Port)) >= min {
			out = append(out, p)
		}
	}
	return out
}
