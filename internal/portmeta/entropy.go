package portmeta

import "math"

// EntropyLevel describes how unpredictable a port's open/close behaviour is.
type EntropyLevel int

const (
	EntropyNone    EntropyLevel = iota // completely predictable
	EntropyLow                         // mostly stable
	EntropyMedium                      // some variability
	EntropyHigh                        // highly unpredictable
)

func (e EntropyLevel) String() string {
	switch e {
	case EntropyNone:
		return "none"
	case EntropyLow:
		return "low"
	case EntropyMedium:
		return "medium"
	case EntropyHigh:
		return "high"
	default:
		return "unknown"
	}
}

// EntropyFor computes an entropy level for a port based on the proportion of
// state-change events relative to total observations. eventCount is the number
// of open/close transitions; sampleCount is the total number of scan samples
// in which the port was observed.
func EntropyFor(eventCount, sampleCount int) EntropyLevel {
	if sampleCount <= 0 || eventCount <= 0 {
		return EntropyNone
	}

	ratio := float64(eventCount) / float64(sampleCount)
	// Normalise to [0,1] and compute a simple Shannon-inspired score.
	// We cap ratio at 1 to avoid log-domain issues.
	if ratio > 1 {
		ratio = 1
	}

	// Use binary entropy H = -p*log2(p) - (1-p)*log2(1-p), scaled to [0,1].
	var h float64
	if ratio > 0 && ratio < 1 {
		h = -ratio*math.Log2(ratio) - (1-ratio)*math.Log2(1-ratio)
	}

	switch {
	case h >= 0.9:
		return EntropyHigh
	case h >= 0.5:
		return EntropyMedium
	case h > 0:
		return EntropyLow
	default:
		return EntropyNone
	}
}
