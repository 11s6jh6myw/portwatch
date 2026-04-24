package portmeta

import (
	"github.com/joshbeard/portwatch/internal/scanner"
)

// CompositeScore represents an aggregated risk/priority score for a port,
// combining multiple signal dimensions into a single comparable value.
type CompositeScore struct {
	// Total is the weighted sum of all component scores, in the range [0, 100].
	Total int

	// Risk contribution from RiskLevel.
	Risk int
	// Severity contribution from SeverityLevel.
	Severity int
	// Exposure contribution from ExposureLevel.
	Exposure int
	// Anomaly contribution from AnomalyLevel.
	Anomaly int
	// Criticality contribution from CriticalityLevel.
	Criticality int
	// Reputation contribution (inverted — low reputation raises the score).
	Reputation int
}

// scoreWeights defines the relative importance of each dimension.
// All weights must sum to 100.
var scoreWeights = struct {
	Risk        int
	Severity    int
	Exposure    int
	Anomaly     int
	Criticality int
	Reputation  int
}{
	Risk:        25,
	Severity:    20,
	Exposure:    20,
	Anomaly:     15,
	Criticality: 10,
	Reputation:  10,
}

// CompositeScoreFor computes an aggregated priority score for the given port.
// Higher scores indicate ports that warrant closer attention. The score is
// normalised to [0, 100].
func CompositeScoreFor(p scanner.PortInfo) CompositeScore {
	risk := int(Score(p))       // 0–3
	sev := int(SeverityFor(p))  // 0–3
	exp := int(ExposureFor(p))  // 0–3
	anom := int(AnomalyFor(p))  // 0–3
	crit := int(CriticalityFor(p)) // 0–3

	// Reputation is inverted: lower reputation → higher contribution.
	rep := 3 - int(ReputationFor(p)) // 0–3

	// Normalise each dimension to a 0–100 scale then apply weight.
	norm := func(raw, weight int) int {
		if raw <= 0 {
			return 0
		}
		return (raw * weight * 100) / (3 * 100)
	}

	riskC := norm(risk, scoreWeights.Risk)
	sevC := norm(sev, scoreWeights.Severity)
	expC := norm(exp, scoreWeights.Exposure)
	anomC := norm(anom, scoreWeights.Anomaly)
	critC := norm(crit, scoreWeights.Criticality)
	repC := norm(rep, scoreWeights.Reputation)

	total := riskC + sevC + expC + anomC + critC + repC
	if total > 100 {
		total = 100
	}

	return CompositeScore{
		Total:       total,
		Risk:        riskC,
		Severity:    sevC,
		Exposure:    expC,
		Anomaly:     anomC,
		Criticality: critC,
		Reputation:  repC,
	}
}

// Priority returns a human-readable priority label derived from the composite
// score total.
func (s CompositeScore) Priority() string {
	switch {
	case s.Total >= 75:
		return "critical"
	case s.Total >= 50:
		return "high"
	case s.Total >= 25:
		return "medium"
	default:
		return "low"
	}
}
