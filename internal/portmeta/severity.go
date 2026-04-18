package portmeta

// Severity represents an alert severity level derived from port metadata.
type Severity int

const (
	SeverityNone Severity = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

func (s Severity) String() string {
	switch s {
	case SeverityLow:
		return "low"
	case SeverityMedium:
		return "medium"
	case SeverityHigh:
		return "high"
	case SeverityCritical:
		return "critical"
	default:
		return "none"
	}
}

// SeverityFor returns a Severity for the given port number by combining
// risk score and category information.
func SeverityFor(port int) Severity {
	meta, ok := Lookup(port)
	if !ok {
		if IsRisky(port) {
			return SeverityMedium
		}
		return SeverityNone
	}

	rl := Score(meta)
	cat := Categorize(port)

	switch rl {
	case RiskHigh:
		if cat == CategoryAdmin || cat == CategoryDatabase {
			return SeverityCritical
		}
		return SeverityHigh
	case RiskMedium:
		return SeverityMedium
	case RiskLow:
		return SeverityLow
	default:
		return SeverityNone
	}
}
