package baseline

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// Violation describes a port that deviates from the saved baseline.
type Violation struct {
	Port   scanner.PortInfo
	Reason string
}

func (v Violation) String() string {
	return fmt.Sprintf("port %d/%s: %s", v.Port.Port, v.Port.Proto, v.Reason)
}

// Check compares current open ports against the baseline snapshot and
// returns any violations (unexpected open or missing expected ports).
func Check(snap *Snapshot, current []scanner.PortInfo) []Violation {
	baseMap := make(map[string]scanner.PortInfo, len(snap.Ports))
	for _, p := range snap.Ports {
		baseMap[key(p)] = p
	}
	currentMap := make(map[string]scanner.PortInfo, len(current))
	for _, p := range current {
		currentMap[key(p)] = p
	}

	var violations []Violation

	// Ports open now but not in baseline.
	for k, p := range currentMap {
		if _, ok := baseMap[k]; !ok {
			violations = append(violations, Violation{
				Port:   p,
				Reason: "not in baseline (unexpected open port)",
			})
		}
	}

	// Ports in baseline but no longer open.
	for k, p := range baseMap {
		if _, ok := currentMap[k]; !ok {
			violations = append(violations, Violation{
				Port:   p,
				Reason: "missing from current scan (expected port closed)",
			})
		}
	}

	return violations
}

func key(p scanner.PortInfo) string {
	return fmt.Sprintf("%d/%s", p.Port, p.Proto)
}
