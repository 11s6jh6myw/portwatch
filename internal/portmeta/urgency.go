package portmeta

import "time"

// UrgencyLevel represents how urgently a port change should be acted upon.
type UrgencyLevel int

const (
	UrgencyNone UrgencyLevel = iota
	UrgencyLow
	UrgencyMedium
	UrgencyHigh
	UrgencyCritical
)

func (u UrgencyLevel) String() string {
	switch u {
	case UrgencyLow:
		return "low"
	case UrgencyMedium:
		return "medium"
	case UrgencyHigh:
		return "high"
	case UrgencyCritical:
		return "critical"
	default:
		return "none"
	}
}

// UrgencyFor computes the urgency level for a port based on its risk, severity,
// exposure, and how recently it was first seen.
func UrgencyFor(port int, firstSeen time.Time) UrgencyLevel {
	risk := Score(port)
	sev := SeverityFor(port)
	exposure := ExposureFor(port)

	score := 0

	switch {
	case risk >= 8:
		score += 3
	case risk >= 5:
		score += 2
	case risk >= 2:
		score += 1
	}

	switch sev {
	case SeverityCritical:
		score += 3
	case SeverityHigh:
		score += 2
	case SeverityMedium:
		score += 1
	}

	switch exposure {
	case ExposureHigh:
		score += 2
	case ExposureMedium:
		score += 1
	}

	// Boost urgency for very recently opened ports (within 5 minutes).
	if !firstSeen.IsZero() && time.Since(firstSeen) < 5*time.Minute {
		score += 1
	}

	switch {
	case score >= 7:
		return UrgencyCritical
	case score >= 5:
		return UrgencyHigh
	case score >= 3:
		return UrgencyMedium
	case score >= 1:
		return UrgencyLow
	default:
		return UrgencyNone
	}
}
