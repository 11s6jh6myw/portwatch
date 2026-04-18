package portmeta

import "github.com/joshbeard/portwatch/internal/scanner"

const trustKey = "trust"

// TrustAnnotator enriches PortInfo with trust level metadata.
type TrustAnnotator struct{}

// NewTrustAnnotator returns a new TrustAnnotator.
func NewTrustAnnotator() *TrustAnnotator {
	return &TrustAnnotator{}
}

// Annotate sets the "trust" metadata field on each port.
func (a *TrustAnnotator) Annotate(ports []scanner.PortInfo) []scanner.PortInfo {
	out := make([]scanner.PortInfo, len(ports))
	for i, p := range ports {
		if p.Meta == nil {
			p.Meta = make(map[string]string)
		}
		p.Meta[trustKey] = TrustFor(uint16(p.Port)).String()
		out[i] = p
	}
	return out
}

// FilterByMinTrust returns only ports whose trust level is >= min.
func FilterByMinTrust(ports []scanner.PortInfo, min TrustLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if TrustFor(uint16(p.Port)) >= min {
			out = append(out, p)
		}
	}
	return out
}
