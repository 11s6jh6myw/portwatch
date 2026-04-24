package portmeta

// AttackSurface represents how much of an attack surface a port presents.
type AttackSurface int

const (
	SurfaceNone    AttackSurface = iota // not exposed or internal only
	SurfaceMinimal                      // limited exposure, well-controlled
	SurfaceModerate                     // standard service, some exposure
	SurfaceSignificant                  // broad exposure or sensitive service
	SurfaceCritical                     // high-value target, maximum exposure
)

// String returns a human-readable label for the attack surface level.
func (a AttackSurface) String() string {
	switch a {
	case SurfaceNone:
		return "none"
	case SurfaceMinimal:
		return "minimal"
	case SurfaceModerate:
		return "moderate"
	case SurfaceSignificant:
		return "significant"
	case SurfaceCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// AttackSurfaceFor derives an attack surface rating for the given port number
// by combining visibility, risk, and exposure signals.
func AttackSurfaceFor(port int) AttackSurface {
	vis := VisibilityFor(port)
	risk := Score(port)
	exp := ExposureFor(port)

	score := 0

	switch vis {
	case VisibilityPublic:
		score += 3
	case VisibilityRestricted:
		score += 2
	case VisibilityInternal:
		score += 1
	}

	switch risk {
	case RiskHigh:
		score += 3
	case RiskMedium:
		score += 2
	case RiskLow:
		score += 1
	}

	switch exp {
	case ExposureHigh:
		score += 3
	case ExposureMedium:
		score += 2
	case ExposureLow:
		score += 1
	}

	switch {
	case score >= 8:
		return SurfaceCritical
	case score >= 6:
		return SurfaceSignificant
	case score >= 4:
		return SurfaceModerate
	case score >= 2:
		return SurfaceMinimal
	default:
		return SurfaceNone
	}
}
