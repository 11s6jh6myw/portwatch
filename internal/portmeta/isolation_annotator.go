package portmeta

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// IsolationAnnotator enriches PortInfo metadata with isolation level.
type IsolationAnnotator struct {
	peers     []scanner.PortInfo
	firstSeen map[uint16]time.Time
}

// NewIsolationAnnotator returns an annotator that uses peers as the context
// population and firstSeen to determine port age.
func NewIsolationAnnotator(peers []scanner.PortInfo, firstSeen map[uint16]time.Time) *IsolationAnnotator {
	if firstSeen == nil {
		firstSeen = make(map[uint16]time.Time)
	}
	return &IsolationAnnotator{peers: peers, firstSeen: firstSeen}
}

// Annotate sets the "isolation" metadata key on each port.
func (a *IsolationAnnotator) Annotate(ports []scanner.PortInfo) []scanner.PortInfo {
	out := make([]scanner.PortInfo, len(ports))
	for i, p := range ports {
		fs := a.firstSeen[p.Port]
		level := IsolationFor(p, a.peers, fs)
		if p.Meta == nil {
			p.Meta = make(map[string]string)
		}
		p.Meta["isolation"] = level.String()
		out[i] = p
	}
	return out
}

// FilterByMinIsolation returns only ports whose isolation level is >= min.
func FilterByMinIsolation(ports []scanner.PortInfo, min IsolationLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		val, ok := p.Meta["isolation"]
		if !ok {
			out = append(out, p)
			continue
		}
		var level IsolationLevel
		switch val {
		case "high":
			level = IsolationHigh
		case "medium":
			level = IsolationMedium
		case "low":
			level = IsolationLow
		default:
			level = IsolationNone
		}
		if level >= min {
			out = append(out, p)
		}
	}
	return out
}
