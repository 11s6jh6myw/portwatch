package portmeta

import "github.com/joshbeard/portwatch/internal/scanner"

// DependencyAnnotator enriches ports with dependency metadata.
type DependencyAnnotator struct{}

// NewDependencyAnnotator returns a new DependencyAnnotator.
func NewDependencyAnnotator() *DependencyAnnotator {
	return &DependencyAnnotator{}
}

// Annotate adds dependency metadata to each port in the slice.
// It sets "dep.count" and "dep.reasons" keys in port metadata.
func (a *DependencyAnnotator) Annotate(ports []scanner.PortInfo) []scanner.PortInfo {
	for i, p := range ports {
		deps := DependenciesFor(p.Port)
		if len(deps) == 0 {
			continue
		}
		if ports[i].Meta == nil {
			ports[i].Meta = make(map[string]string)
		}
		ports[i].Meta["dep.count"] = itoa(len(deps))
		var reasons string
		for j, d := range deps {
			if j > 0 {
				reasons += "; "
			}
			reasons += d.Reason
		}
		ports[i].Meta["dep.reasons"] = reasons
	}
	return ports
}

// FilterByHasDependency returns only ports that have at least one known dependency.
func FilterByHasDependency(ports []scanner.PortInfo) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if len(DependenciesFor(p.Port)) > 0 {
			out = append(out, p)
		}
	}
	return out
}
