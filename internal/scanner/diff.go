package scanner

import "fmt"

// ChangeType represents the type of port change detected.
type ChangeType string

const (
	ChangeOpened ChangeType = "OPENED"
	ChangeClosed ChangeType = "CLOSED"
)

// PortChange describes a single detected change in port state.
type PortChange struct {
	Type ChangeType
	Port PortInfo
}

// String returns a human-readable description of the change.
func (c PortChange) String() string {
	return fmt.Sprintf("[%s] %s", c.Type, c.Port)
}

// Diff compares two port snapshots and returns the list of changes.
// previous is the last known state; current is the newly scanned state.
func Diff(previous, current []PortInfo) []PortChange {
	var changes []PortChange

	prevMap := toMap(previous)
	currMap := toMap(current)

	// Detect newly opened ports.
	for key, info := range currMap {
		if _, exists := prevMap[key]; !exists {
			changes = append(changes, PortChange{Type: ChangeOpened, Port: info})
		}
	}

	// Detect closed ports.
	for key, info := range prevMap {
		if _, exists := currMap[key]; !exists {
			changes = append(changes, PortChange{Type: ChangeClosed, Port: info})
		}
	}

	return changes
}

// toMap converts a slice of PortInfo into a map keyed by "protocol:address:port".
func toMap(ports []PortInfo) map[string]PortInfo {
	m := make(map[string]PortInfo, len(ports))
	for _, p := range ports {
		key := fmt.Sprintf("%s:%s:%d", p.Protocol, p.Address, p.Port)
		m[key] = p
	}
	return m
}
