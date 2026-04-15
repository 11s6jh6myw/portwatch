package scanner

// DiffResult holds the changes between two port scans.
type DiffResult struct {
	Opened []PortInfo
	Closed []PortInfo
}

// HasChanges returns true if there are any opened or closed ports.
func (d DiffResult) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

// Diff computes the difference between a previous and current port scan.
func Diff(previous, current []PortInfo) DiffResult {
	prev := toMap(previous)
	curr := toMap(current)

	var result DiffResult

	for port, info := range curr {
		if _, existed := prev[port]; !existed {
			result.Opened = append(result.Opened, info)
		}
	}

	for port, info := range prev {
		if _, exists := curr[port]; !exists {
			result.Closed = append(result.Closed, info)
		}
	}

	return result
}

func toMap(ports []PortInfo) map[int]PortInfo {
	m := make(map[int]PortInfo, len(ports))
	for _, p := range ports {
		m[p.Port] = p
	}
	return m
}
