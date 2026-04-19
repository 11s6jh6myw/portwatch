package portmeta

import "time"

// VelocityLevel describes how rapidly a port's state is changing.
type VelocityLevel int

const (
	VelocityNone VelocityLevel = iota
	VelocitySlow
	VelocityModerate
	VelocityFast
	VelocityRapid
)

func (v VelocityLevel) String() string {
	switch v {
	case VelocitySlow:
		return "slow"
	case VelocityModerate:
		return "moderate"
	case VelocityFast:
		return "fast"
	case VelocityRapid:
		return "rapid"
	default:
		return "none"
	}
}

// VelocityFor computes the rate of state changes for a port given a list of
// event timestamps within the observation window.
func VelocityFor(events []time.Time, window time.Duration) VelocityLevel {
	if len(events) == 0 || window <= 0 {
		return VelocityNone
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
		return VelocityNone
	case count <= 2:
		return VelocitySlow
	case count <= 5:
		return VelocityModerate
	case count <= 10:
		return VelocityFast
	default:
		return VelocityRapid
	}
}
