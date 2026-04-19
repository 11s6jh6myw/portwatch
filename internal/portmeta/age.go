package portmeta

import "time"

// AgeClass describes how long a port has been observed open.
type AgeClass int

const (
	AgeUnknown AgeClass = iota
	AgeFresh             // seen for less than 1 hour
	AgeShortLived       // 1 hour – 24 hours
	AgeMature           // 1 day – 7 days
	AgeEstablished      // more than 7 days
)

func (a AgeClass) String() string {
	switch a {
	case AgeFresh:
		return "fresh"
	case AgeShortLived:
		return "short-lived"
	case AgeMature:
		return "mature"
	case AgeEstablished:
		return "established"
	default:
		return "unknown"
	}
}

// ClassifyAge returns the AgeClass for a port first seen at firstSeen,
// evaluated relative to now.
func ClassifyAge(firstSeen, now time.Time) AgeClass {
	if firstSeen.IsZero() {
		return AgeUnknown
	}
	age := now.Sub(firstSeen)
	switch {
	case age < time.Hour:
		return AgeFresh
	case age < 24*time.Hour:
		return AgeShortLived
	case age < 7*24*time.Hour:
		return AgeMature
	default:
		return AgeEstablished
	}
}
