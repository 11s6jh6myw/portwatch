package portmeta

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

const signalKey = "signal"

// NewSignalAnnotator returns an annotator that adds a "signal" metadata
// entry to each port reflecting its computed SignalStrength.
func NewSignalAnnotator(now time.Time) func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			s := SignalFor(p, now)
			if p.Meta == nil {
				p.Meta = make(map[string]string)
			}
			p.Meta[signalKey] = s.String()
			out[i] = p
		}
		return out
	}
}

// FilterByMinSignal returns only ports whose stored signal metadata is
// at or above the given minimum SignalStrength. Ports without metadata
// are passed through unchanged.
func FilterByMinSignal(ports []scanner.PortInfo, min SignalStrength) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if p.Meta == nil {
			out = append(out, p)
			continue
		}
		v, ok := p.Meta[signalKey]
		if !ok {
			out = append(out, p)
			continue
		}
		if parseSignal(v) >= min {
			out = append(out, p)
		}
	}
	return out
}

func parseSignal(s string) SignalStrength {
	switch s {
	case "weak":
		return SignalWeak
	case "moderate":
		return SignalModerate
	case "strong":
		return SignalStrong
	case "critical":
		return SignalCritical
	default:
		return SignalNone
	}
}
