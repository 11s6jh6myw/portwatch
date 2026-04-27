package portmeta

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// SignalStrength represents how strongly a port's behaviour indicates
// something noteworthy (e.g. a threat, anomaly, or operational event).
type SignalStrength int

const (
	SignalNone   SignalStrength = iota // no meaningful signal
	SignalWeak                         // minor indicator
	SignalModerate                     // worth investigating
	SignalStrong                       // high-confidence indicator
	SignalCritical                     // immediate action recommended
)

func (s SignalStrength) String() string {
	switch s {
	case SignalWeak:
		return "weak"
	case SignalModerate:
		return "moderate"
	case SignalStrong:
		return "strong"
	case SignalCritical:
		return "critical"
	default:
		return "none"
	}
}

// SignalFor computes a composite SignalStrength for a port based on its
// risk, anomaly, urgency, and how recently it was first seen.
func SignalFor(p scanner.PortInfo, now time.Time) SignalStrength {
	score := 0

	if IsRisky(p.Port) {
		score += 3
	}

	anomaly := AnomalyFor(p)
	score += int(anomaly)

	urgency := UrgencyFor(p, now)
	score += int(urgency)

	if !p.FirstSeen.IsZero() && now.Sub(p.FirstSeen) < 5*time.Minute {
		score += 2
	}

	switch {
	case score >= 9:
		return SignalCritical
	case score >= 6:
		return SignalStrong
	case score >= 3:
		return SignalModerate
	case score >= 1:
		return SignalWeak
	default:
		return SignalNone
	}
}
