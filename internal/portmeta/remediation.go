package portmeta

// RemediationLevel indicates the urgency of remediation action required.
type RemediationLevel int

const (
	RemediationNone RemediationLevel = iota
	RemediationMonitor
	RemediationReview
	RemediationMitigate
	RemediationImmediate
)

func (r RemediationLevel) String() string {
	switch r {
	case RemediationMonitor:
		return "monitor"
	case RemediationReview:
		return "review"
	case RemediationMitigate:
		return "mitigate"
	case RemediationImmediate:
		return "immediate"
	default:
		return "none"
	}
}

// RemediationFor returns the recommended remediation level for a given port
// based on its risk, severity, and anomaly characteristics.
func RemediationFor(port int, firstSeen int64, eventCount int, scanCount int) RemediationLevel {
	risk := Score(port, firstSeen, eventCount, scanCount)
	anomaly := AnomalyFor(port, firstSeen, eventCount, scanCount)
	severity := SeverityFor(port)

	if anomaly >= AnomalyHigh && severity >= SeverityCritical {
		return RemediationImmediate
	}
	if risk >= 80 || (anomaly >= AnomalyHigh && severity >= SeverityHigh) {
		return RemediationMitigate
	}
	if risk >= 50 || anomaly >= AnomalyMedium {
		return RemediationReview
	}
	if risk >= 25 || severity >= SeverityMedium {
		return RemediationMonitor
	}
	return RemediationNone
}

// IsActionable returns true if the remediation level requires human action.
func IsActionable(level RemediationLevel) bool {
	return level >= RemediationReview
}
