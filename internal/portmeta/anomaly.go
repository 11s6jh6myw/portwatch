package portmeta

// AnomalyLevel describes how unusual a port's behaviour appears.
type AnomalyLevel int

const (
	AnomalyNone AnomalyLevel = iota
	AnomalyLow
	AnomalyMedium
	AnomalyHigh
)

func (a AnomalyLevel) String() string {
	switch a {
	case AnomalyLow:
		return "low"
	case AnomalyMedium:
		return "medium"
	case AnomalyHigh:
		return "high"
	default:
		return "none"
	}
}

// AnomalyFor returns an anomaly level for the given port number based on
// a combination of risk, trust, and exposure heuristics.
func AnomalyFor(port int) AnomalyLevel {
	risk := Score(port)
	trust := TrustFor(port)
	exposure := ExposureFor(port)

	score := 0
	if risk >= RiskHigh {
		score += 3
	} else if risk >= RiskMedium {
		score += 2
	} else if risk >= RiskLow {
		score++
	}

	if trust <= TrustLow {
		score += 2
	} else if trust <= TrustMedium {
		score++
	}

	if exposure >= ExposureHigh {
		score += 2
	} else if exposure >= ExposureMedium {
		score++
	}

	switch {
	case score >= 6:
		return AnomalyHigh
	case score >= 4:
		return AnomalyMedium
	case score >= 2:
		return AnomalyLow
	default:
		return AnomalyNone
	}
}
