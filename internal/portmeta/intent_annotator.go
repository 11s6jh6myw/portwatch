package portmeta

import "github.com/user/portwatch/internal/scanner"

// NewIntentAnnotator returns an annotator that adds intent metadata to ports.
func NewIntentAnnotator() func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			level := IntentFor(p.Port)
			if p.Meta == nil {
				p.Meta = make(map[string]string)
			}
			p.Meta["intent"] = level.String()
			out[i] = p
		}
		return out
	}
}

// FilterByIntent returns only ports whose intent matches the given level.
func FilterByIntent(ports []scanner.PortInfo, intent IntentLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if IntentFor(p.Port) == intent {
			out = append(out, p)
		}
	}
	return out
}
