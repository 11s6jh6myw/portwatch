package portmeta

import "github.com/iamcalledrob/portwatch/internal/scanner"

// ImpactLevel describes the potential blast radius if a port is compromised or unexpectedly exposed.
type ImpactLevel int

const (
	ImpactNone ImpactLevel = iota
	ImpactLow
	ImpactMedium
	ImpactHigh
	ImpactCritical
)

func (l ImpactLevel) String() string {
	switch l {
	case ImpactLow:
		return "low"
	case ImpactMedium:
		return "medium"
	case ImpactHigh:
		return "high"
	case ImpactCritical:
		return "critical"
	default:
		return "none"
	}
}

// ImpactFor returns the estimated impact level for the given port.
func ImpactFor(p scanner.PortInfo) ImpactLevel {
	meta, ok := Lookup(p.Port)
	if !ok {
		return ImpactNone
	}

	// Combine risk and criticality signals.
	crit := CriticalityFor(p)
	risk := Score(p)

	switch {
	case crit == CriticalityCritical || (risk >= 9 && meta.Known):
		return ImpactCritical
	case crit == CriticalityHigh || risk >= 7:
		return ImpactHigh
	case crit == CriticalityMedium || risk >= 4:
		return ImpactMedium
	case meta.Known:
		return ImpactLow
	default:
		return ImpactNone
	}
}

// IsHighImpact returns true when the port's impact is High or Critical.
func IsHighImpact(p scanner.PortInfo) bool {
	l := ImpactFor(p)
	return l == ImpactHigh || l == ImpactCritical
}
