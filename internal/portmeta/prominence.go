package portmeta

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// ProminenceLevel indicates how notable or significant a port is within the
// observed landscape, combining visibility, prevalence, and activity signals.
type ProminenceLevel int

const (
	ProminenceNone ProminenceLevel = iota
	ProminenceLow
	ProminenceMedium
	ProminenceHigh
	ProminenceCritical
)

func (p ProminenceLevel) String() string {
	switch p {
	case ProminenceLow:
		return "low"
	case ProminenceMedium:
		return "medium"
	case ProminenceHigh:
		return "high"
	case ProminenceCritical:
		return "critical"
	default:
		return "none"
	}
}

// ProminenceFor computes a prominence level for the given port based on its
// metadata. Ports seen frequently, flagged as well-known, or recently active
// score higher.
func ProminenceFor(p scanner.PortInfo) ProminenceLevel {
	if p.Meta == nil {
		return ProminenceNone
	}

	score := 0

	if IsPrevalent(p) {
		score += 2
	}

	if IsCritical(p) {
		score += 2
	}

	if IsHighImpact(p) {
		score++
	}

	if lastSeen, ok := p.Meta["last_seen"]; ok {
		if t, err := time.Parse(time.RFC3339, lastSeen); err == nil {
			if time.Since(t) < 5*time.Minute {
				score++
			}
		}
	}

	switch {
	case score >= 5:
		return ProminenceCritical
	case score >= 3:
		return ProminenceHigh
	case score >= 2:
		return ProminenceMedium
	case score >= 1:
		return ProminenceLow
	default:
		return ProminenceNone
	}
}
