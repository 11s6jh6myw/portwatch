package portmeta

import "time"

// FreshnessLevel describes how recently a port was first seen.
type FreshnessLevel int

const (
	FreshnessUnknown  FreshnessLevel = iota
	FreshnessStale                   // > 30 days
	FreshnessMature                  // 7–30 days
	FreshnessRecent                  // 1–7 days
	FreshnessNew                     // < 24 hours
)

func (f FreshnessLevel) String() string {
	switch f {
	case FreshnessNew:
		return "new"
	case FreshnessRecent:
		return "recent"
	case FreshnessMature:
		return "mature"
	case FreshnessStale:
		return "stale"
	default:
		return "unknown"
	}
}

// ClassifyFreshness returns a FreshnessLevel based on how long ago firstSeen was.
// A zero firstSeen returns FreshnessUnknown.
func ClassifyFreshness(firstSeen time.Time) FreshnessLevel {
	if firstSeen.IsZero() {
		return FreshnessUnknown
	}
	age := time.Since(firstSeen)
	switch {
	case age < 24*time.Hour:
		return FreshnessNew
	case age < 7*24*time.Hour:
		return FreshnessRecent
	case age < 30*24*time.Hour:
		return FreshnessMature
	default:
		return FreshnessStale
	}
}
