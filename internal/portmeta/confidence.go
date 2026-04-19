package portmeta

import "time"

// ConfidenceLevel represents how confident we are in a port classification.
type ConfidenceLevel int

const (
	ConfidenceLow ConfidenceLevel = iota
	ConfidenceMedium
	ConfidenceHigh
	ConfidenceCertain
)

func (c ConfidenceLevel) String() string {
	switch c {
	case ConfidenceLow:
		return "low"
	case ConfidenceMedium:
		return "medium"
	case ConfidenceHigh:
		return "high"
	case ConfidenceCertain:
		return "certain"
	default:
		return "unknown"
	}
}

// ConfidenceFor returns a confidence level for a port based on how long it has
// been observed and how frequently it changes.
func ConfidenceFor(firstSeen time.Time, observations int, freq ChangeFreq) ConfidenceLevel {
	if firstSeen.IsZero() || observations == 0 {
		return ConfidenceLow
	}

	age := time.Since(firstSeen)

	switch {
	case freq == ChangeFreqVolatile:
		return ConfidenceLow
	case observations >= 50 && age >= 7*24*time.Hour && freq <= ChangeFreqOccasional:
		return ConfidenceCertain
	case observations >= 20 && age >= 24*time.Hour:
		return ConfidenceHigh
	case observations >= 5 && age >= time.Hour:
		return ConfidenceMedium
	default:
		return ConfidenceLow
	}
}
