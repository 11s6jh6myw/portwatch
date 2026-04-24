package portmeta

import (
	"strconv"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

const (
	noiseKey       = "noise.level"
	noiseCountKey  = "noise.event_count"
)

// NewNoiseAnnotator returns an annotator that records the noise level and
// event count for each port using a fixed observation window.
func NewNoiseAnnotator(events map[int][]time.Time, window time.Duration) func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			evs := events[p.Port]
			lvl := NoiseFor(p, evs, window)

			cutoff := time.Now().Add(-window)
			count := 0
			for _, t := range evs {
				if t.After(cutoff) {
					count++
				}
			}

			if p.Metadata == nil {
				p.Metadata = make(map[string]string)
			}
			p.Metadata[noiseKey] = lvl.String()
			p.Metadata[noiseCountKey] = strconv.Itoa(count)
			out[i] = p
		}
		return out
	}
}

// FilterByMaxNoise returns only ports whose recorded noise level is at or
// below the given maximum. Ports without noise metadata are passed through.
func FilterByMaxNoise(ports []scanner.PortInfo, max NoiseLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		raw, ok := p.Metadata[noiseKey]
		if !ok {
			out = append(out, p)
			continue
		}
		lvl := parseNoise(raw)
		if lvl <= max {
			out = append(out, p)
		}
	}
	return out
}

func parseNoise(s string) NoiseLevel {
	switch s {
	case "minimal":
		return NoiseMinimal
	case "moderate":
		return NoiseModerate
	case "high":
		return NoiseHigh
	default:
		return NoiseNone
	}
}
