package portmeta

import "github.com/user/portwatch/internal/scanner"

// ScoreLevel represents a normalised 0–100 score bucket.
type ScoreLevel int

const (
	ScoreNegligible ScoreLevel = iota // 0–19
	ScoreLow                          // 20–39
	ScoreMedium                       // 40–59
	ScoreHigh                         // 60–79
	ScoreCritical                     // 80–100
)

func (s ScoreLevel) String() string {
	switch s {
	case ScoreNegligible:
		return "negligible"
	case ScoreLow:
		return "low"
	case ScoreMedium:
		return "medium"
	case ScoreHigh:
		return "high"
	case ScoreCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// BucketScore maps a raw 0–100 integer into a ScoreLevel.
func BucketScore(score int) ScoreLevel {
	switch {
	case score >= 80:
		return ScoreCritical
	case score >= 60:
		return ScoreHigh
	case score >= 40:
		return ScoreMedium
	case score >= 20:
		return ScoreLow
	default:
		return ScoreNegligible
	}
}

// NormalisedScore returns BucketScore(CompositeScoreFor(p)).
func NormalisedScore(p scanner.PortInfo) ScoreLevel {
	return BucketScore(CompositeScoreFor(p))
}
