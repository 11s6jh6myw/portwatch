package portmeta

import "time"

// VolatilityLevel describes how unpredictably a port's state changes over time.
type VolatilityLevel int

const (
	VolatilityNone     VolatilityLevel = iota // never observed changing
	VolatilityLow                             // rarely changes
	VolatilityModerate                        // changes occasionally
	VolatilityHigh                            // changes frequently and unpredictably
)

func (v VolatilityLevel) String() string {
	switch v {
	case VolatilityNone:
		return "none"
	case VolatilityLow:
		return "low"
	case VolatilityModerate:
		return "moderate"
	case VolatilityHigh:
		return "high"
	default:
		return "unknown"
	}
}

// VolatilityFor derives a volatility level from the number of state-change
// events observed within the provided window duration.
//
// events is the slice of timestamps at which the port changed state (opened or
// closed). Only events within [now-window, now] are considered.
func VolatilityFor(events []time.Time, window time.Duration) VolatilityLevel {
	if len(events) == 0 || window <= 0 {
		return VolatilityNone
	}

	cutoff := time.Now().Add(-window)
	count := 0
	for _, e := range events {
		if e.After(cutoff) {
			count++
		}
	}

	switch {
	case count == 0:
		return VolatilityNone
	case count <= 2:
		return VolatilityLow
	case count <= 6:
		return VolatilityModerate
	default:
		return VolatilityHigh
	}
}
