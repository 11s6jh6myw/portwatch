package portmeta

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// NoiseLevel describes how much irrelevant or transient activity a port exhibits.
type NoiseLevel int

const (
	NoiseNone    NoiseLevel = iota // no transient activity
	NoiseMinimal                   // rare, isolated events
	NoiseModerate                  // occasional spurious events
	NoiseHigh                      // frequent transient open/close cycles
)

func (n NoiseLevel) String() string {
	switch n {
	case NoiseNone:
		return "none"
	case NoiseMinimal:
		return "minimal"
	case NoiseModerate:
		return "moderate"
	case NoiseHigh:
		return "high"
	default:
		return "unknown"
	}
}

// NoiseFor estimates the noise level for a port based on event frequency and
// the ratio of open/close transitions within the observation window.
func NoiseFor(p scanner.PortInfo, events []time.Time, window time.Duration) NoiseLevel {
	if window <= 0 || len(events) == 0 {
		return NoiseNone
	}

	cutoff := time.Now().Add(-window)
	count := 0
	for _, t := range events {
		if t.After(cutoff) {
			count++
		}
	}

	// Combine port-level risk with raw event count to estimate noise.
	multiplier := 1.0
	if IsRisky(p.Port) {
		multiplier = 0.75 // risky ports are expected to be noisier; lower threshold
	}

	adjusted := float64(count) * multiplier
	switch {
	case adjusted >= 20:
		return NoiseHigh
	case adjusted >= 8:
		return NoiseModerate
	case adjusted >= 2:
		return NoiseMinimal
	default:
		return NoiseNone
	}
}
