package portmeta

import "github.com/user/portwatch/internal/scanner"

const surfaceKey = "attack_surface"

// NewSurfaceAnnotator returns a function that annotates each PortInfo with
// its computed attack surface level stored under the "attack_surface" metadata key.
func NewSurfaceAnnotator() func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			if p.Meta == nil {
				p.Meta = make(map[string]string)
			}
			p.Meta[surfaceKey] = AttackSurfaceFor(p.Port).String()
			out[i] = p
		}
		return out
	}
}

// FilterByMinSurface returns only the ports whose attack surface is at or
// above the given minimum level.
func FilterByMinSurface(ports []scanner.PortInfo, min AttackSurface) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if p.Meta == nil {
			out = append(out, p)
			continue
		}
		level := parseSurface(p.Meta[surfaceKey])
		if level >= min {
			out = append(out, p)
		}
	}
	return out
}

func parseSurface(s string) AttackSurface {
	switch s {
	case "critical":
		return SurfaceCritical
	case "significant":
		return SurfaceSignificant
	case "moderate":
		return SurfaceModerate
	case "minimal":
		return SurfaceMinimal
	default:
		return SurfaceNone
	}
}
