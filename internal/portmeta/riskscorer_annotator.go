package portmeta

import "github.com/user/portwatch/internal/scanner"

// Annotated wraps a PortInfo with its computed risk level.
type Annotated struct {
	Port  scanner.PortInfo
	Risk  RiskLevel
	Label string
}

// Annotate enriches a slice of PortInfo with risk scores and labels.
func Annotate(ports []scanner.PortInfo) []Annotated {
	out := make([]Annotated, 0, len(ports))
	for _, p := range ports {
		meta, _ := Lookup(p.Port)
		label := meta
		if label == "" {
			label = "unknown"
		}
		out = append(out, Annotated{
			Port:  p,
			Risk:  Score(p.Port),
			Label: label,
		})
	}
	return out
}

// FilterByMinRisk returns only those Annotated entries at or above minRisk.
func FilterByMinRisk(annotated []Annotated, minRisk RiskLevel) []Annotated {
	var out []Annotated
	for _, a := range annotated {
		if a.Risk >= minRisk {
			out = append(out, a)
		}
	}
	return out
}
