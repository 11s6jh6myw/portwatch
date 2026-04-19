package portmeta

import "time"

// StabilityLevel describes how stable a port's presence has been over time.
type StabilityLevel int

const (
	StabilityUnknown StabilityLevel = iota
	StabilityUnstable
	StabilityVariable
	StabilityStable
	StabilityLocked
)

func (s StabilityLevel) String() string {
	switch s {
	case StabilityUnstable:
		return "unstable"
	case StabilityVariable:
		return "variable"
	case StabilityStable:
		return "stable"
	case StabilityLocked:
		return "locked"
	default:
		return "unknown"
	}
}

// ClassifyStability returns a StabilityLevel based on how long the port has
// been continuously open and how many state changes it has had.
func ClassifyStability(firstSeen time.Time, changeCount int) StabilityLevel {
	if firstSeen.IsZero() {
		return StabilityUnknown
	}
	age := time.Since(firstSeen)
	switch {
	case changeCount >= 10:
		return StabilityUnstable
	case changeCount >= 4:
		return StabilityVariable
	case age >= 30*24*time.Hour && changeCount <= 1:
		return StabilityLocked
	default:
		return StabilityStable
	}
}
