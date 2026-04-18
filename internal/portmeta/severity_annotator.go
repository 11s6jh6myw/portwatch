package portmeta

import "github.com/netwatch/portwatch/internal/scanner"

// SeverityAnnotator attaches a severity tag to each PortInfo.
type SeverityAnnotator struct{}

// NewSeverityAnnotator returns a new SeverityAnnotator.
func NewSeverityAnnotator() *SeverityAnnotator {
	return &SeverityAnnotator{}
}

// Annotate returns a copy of ports with a "severity" tag added.
func (a *SeverityAnnotator) Annotate(ports []scanner.PortInfo) []scanner.PortInfo {
	out := make([]scanner.PortInfo, len(ports))
	for i, p := range ports {
		sev := SeverityFor(p.Port)
		if p.Tags == nil {
			p.Tags = map[string]string{}
		} else {
			tags := make(map[string]string, len(p.Tags)+1)
			for k, v := range p.Tags {
				tags[k] = v
			}
			p.Tags = tags
		}
		p.Tags["severity"] = sev.String()
		out[i] = p
	}
	return out
}

// FilterByMinSeverity returns only ports whose severity is >= min.
func FilterByMinSeverity(ports []scanner.PortInfo, min Severity) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if SeverityFor(p.Port) >= min {
			out = append(out, p)
		}
	}
	return out
}
