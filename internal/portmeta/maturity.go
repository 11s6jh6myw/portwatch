package portmeta

import "time"

// MaturityLevel describes how mature a port's presence is considered.
type MaturityLevel int

const (
	MaturityUnknown MaturityLevel = iota
	MaturityEmerging
	MaturityDeveloping
	MaturityEstablished
	MaturityMature
)

func (m MaturityLevel) String() string {
	switch m {
	case MaturityEmerging:
		return "emerging"
	case MaturityDeveloping:
		return "developing"
	case MaturityEstablished:
		return "established"
	case MaturityMature:
		return "mature"
	default:
		return "unknown"
	}
}

// MaturityFor returns a MaturityLevel based on first-seen time and event count.
func MaturityFor(firstSeen time.Time, eventCount int) MaturityLevel {
	if firstSeen.IsZero() {
		return MaturityUnknown
	}
	age := time.Since(firstSeen)
	switch {
	case age < 24*time.Hour && eventCount < 5:
		return MaturityEmerging
	case age < 7*24*time.Hour && eventCount < 20:
		return MaturityDeveloping
	case age < 30*24*time.Hour:
		return MaturityEstablished
	default:
		return MaturityMature
	}
}
