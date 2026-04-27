package portmeta

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// PressureLevel indicates how much cumulative scanning or connection
// pressure a port has experienced over a recent observation window.
type PressureLevel int

const (
	PressureNone     PressureLevel = iota // no recorded activity
	PressureLow                           // fewer than 5 events
	PressureModerate                      // 5–19 events
	PressureHigh                          // 20–49 events
	PressureCritical                      // 50+ events
)

func (p PressureLevel) String() string {
	switch p {
	case PressureLow:
		return "low"
	case PressureModerate:
		return "moderate"
	case PressureHigh:
		return "high"
	case PressureCritical:
		return "critical"
	default:
		return "none"
	}
}

// PressureFor computes the pressure level for a port based on the number
// of events recorded within the given window duration. A zero window
// defaults to 1 hour.
func PressureFor(port scanner.PortInfo, events []time.Time, window time.Duration) PressureLevel {
	if window <= 0 {
		window = time.Hour
	}
	cutoff := time.Now().Add(-window)
	count := 0
	for _, t := range events {
		if t.After(cutoff) {
			count++
		}
	}
	switch {
	case count == 0:
		return PressureNone
	case count < 5:
		return PressureLow
	case count < 20:
		return PressureModerate
	case count < 50:
		return PressureHigh
	default:
		return PressureCritical
	}
}
