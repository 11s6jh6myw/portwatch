package portmeta

import "github.com/joeshaw/portwatch/internal/scanner"

// NewVisibilityAnnotator returns a function that annotates each PortInfo
// with a "visibility" metadata key.
func NewVisibilityAnnotator() func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			v := VisibilityFor(p.Port)
			if p.Meta == nil {
				p.Meta = make(map[string]string)
			}
			p.Meta["visibility"] = v.String()
			out[i] = p
		}
		return out
	}
}

// FilterByMinVisibility returns only ports whose visibility is >= min.
func FilterByMinVisibility(ports []scanner.PortInfo, min VisibilityLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if VisibilityFor(p.Port) >= min {
			out = append(out, p)
		}
	}
	return out
}
