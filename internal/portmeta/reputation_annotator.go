package portmeta

import "github.com/iamralch/portwatch/internal/scanner"

// ReputationAnnotator enriches PortInfo with reputation metadata.
type ReputationAnnotator struct{}

// NewReputationAnnotator returns a new ReputationAnnotator.
func NewReputationAnnotator() *ReputationAnnotator {
	return &ReputationAnnotator{}
}

// Annotate sets "reputation" and "reputable" metadata keys on each port.
func (a *ReputationAnnotator) Annotate(ports []scanner.PortInfo) []scanner.PortInfo {
	out := make([]scanner.PortInfo, len(ports))
	for i, p := range ports {
		r := ReputationFor(uint16(p.Port))
		if p.Meta == nil {
			p.Meta = make(map[string]string)
		}
		p.Meta["reputation"] = r.String()
		if IsReputable(uint16(p.Port)) {
			p.Meta["reputable"] = "true"
		} else {
			p.Meta["reputable"] = "false"
		}
		out[i] = p
	}
	return out
}

// FilterByMinReputation returns only ports whose reputation is >= min.
func FilterByMinReputation(ports []scanner.PortInfo, min ReputationLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if ReputationFor(uint16(p.Port)) >= min {
			out = append(out, p)
		}
	}
	return out
}
