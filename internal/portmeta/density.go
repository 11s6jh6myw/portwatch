package portmeta

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// DensityLevel represents how densely a port has been observed across scans.
type DensityLevel int

const (
	DensityNone DensityLevel = iota
	DensitySparse
	DensityModerate
	DensityDense
	DensitySaturated
)

func (d DensityLevel) String() string {
	switch d {
	case DensitySparse:
		return "sparse"
	case DensityModerate:
		return "moderate"
	case DensityDense:
		return "dense"
	case DensitySaturated:
		return "saturated"
	default:
		return "none"
	}
}

// DensityFor computes how frequently a port appeared across a set of scan
// observations within the given time window. The ratio of appearances to total
// possible scan slots drives the classification.
func DensityFor(p scanner.PortInfo, events []time.Time, window time.Duration, totalScans int) DensityLevel {
	if totalScans <= 0 || len(events) == 0 || window <= 0 {
		return DensityNone
	}

	cutoff := time.Now().Add(-window)
	count := 0
	for _, t := range events {
		if t.After(cutoff) {
			count++
		}
	}

	if count == 0 {
		return DensityNone
	}

	ratio := float64(count) / float64(totalScans)
	switch {
	case ratio >= 0.9:
		return DensitySaturated
	case ratio >= 0.65:
		return DensityDense
	case ratio >= 0.35:
		return DensityModerate
	default:
		return DensitySparse
	}
}
