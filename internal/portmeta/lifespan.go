package portmeta

import (
	"time"
)

// LifespanLevel represents how long a port has been continuously observed open.
type LifespanLevel int

const (
	LifespanUnknown   LifespanLevel = iota
	LifespanEphemeral               // < 1 minute
	LifespanShort                   // 1m – 1h
	LifespanMedium                  // 1h – 24h
	LifespanLong                    // 24h – 7d
	LifespanPermanent               // > 7d
)

func (l LifespanLevel) String() string {
	switch l {
	case LifespanEphemeral:
		return "ephemeral"
	case LifespanShort:
		return "short"
	case LifespanMedium:
		return "medium"
	case LifespanLong:
		return "long"
	case LifespanPermanent:
		return "permanent"
	default:
		return "unknown"
	}
}

// LifespanFor returns the lifespan level based on the first-seen time.
// A zero firstSeen returns LifespanUnknown.
func LifespanFor(firstSeen time.Time) LifespanLevel {
	if firstSeen.IsZero() {
		return LifespanUnknown
	}
	age := time.Since(firstSeen)
	switch {
	case age < time.Minute:
		return LifespanEphemeral
	case age < time.Hour:
		return LifespanShort
	case age < 24*time.Hour:
		return LifespanMedium
	case age < 7*24*time.Hour:
		return LifespanLong
	default:
		return LifespanPermanent
	}
}

// IsLongLived returns true if the port has been open for at least 24 hours.
func IsLongLived(firstSeen time.Time) bool {
	l := LifespanFor(firstSeen)
	return l == LifespanLong || l == LifespanPermanent
}
