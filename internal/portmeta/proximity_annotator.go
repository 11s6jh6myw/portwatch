package portmeta

import (
	"strconv"

	"github.com/user/portwatch/internal/scanner"
)

const proximityKey = "proximity"

// NewProximityAnnotator returns an annotator that enriches each PortInfo with
// its proximity level relative to all other ports in the same scan result.
func NewProximityAnnotator(all []scanner.PortInfo) func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			level := ProximityFor(p, all)
			if p.Meta == nil {
				p.Meta = make(map[string]string)
			}
			p.Meta[proximityKey] = level.String()
			p.Meta[proximityKey+"_score"] = strconv.Itoa(int(level))
			out[i] = p
		}
		return out
	}
}

// FilterByMinProximity returns only ports whose annotated proximity level is
// at or above the given minimum. Ports without proximity metadata pass through.
func FilterByMinProximity(ports []scanner.PortInfo, min ProximityLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if p.Meta == nil {
			out = append(out, p)
			continue
		}
		v, ok := p.Meta[proximityKey+"_score"]
		if !ok {
			out = append(out, p)
			continue
		}
		score, err := strconv.Atoi(v)
		if err != nil {
			out = append(out, p)
			continue
		}
		if ProximityLevel(score) >= min {
			out = append(out, p)
		}
	}
	return out
}
