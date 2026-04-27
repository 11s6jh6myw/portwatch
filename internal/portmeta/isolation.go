package portmeta

import "github.com/iamcalledrob/portwatch/internal/scanner"

// IsolationLevel describes how isolated a port is from other services.
type IsolationLevel int

const (
	IsolationNone    IsolationLevel = iota // co-located with many known services
	IsolationLow                           // shares address space with a few services
	IsolationMedium                        // loosely grouped
	IsolationHigh                          // no known co-located services
)

func (l IsolationLevel) String() string {
	switch l {
	case IsolationNone:
		return "none"
	case IsolationLow:
		return "low"
	case IsolationMedium:
		return "medium"
	case IsolationHigh:
		return "high"
	default:
		return "unknown"
	}
}

// IsolationFor returns the IsolationLevel for a port given the full set of
// currently open ports. A port is considered more isolated when fewer
// well-known neighbours share the same host.
func IsolationFor(p scanner.PortInfo, peers []scanner.PortInfo) IsolationLevel {
	if len(peers) == 0 {
		return IsolationHigh
	}

	knownNeighbours := 0
	for _, peer := range peers {
		if peer.Port == p.Port {
			continue
		}
		if _, ok := wellKnownPorts[peer.Port]; ok {
			knownNeighbours++
		}
	}

	switch {
	case knownNeighbours == 0:
		return IsolationHigh
	case knownNeighbours <= 2:
		return IsolationMedium
	case knownNeighbours <= 5:
		return IsolationLow
	default:
		return IsolationNone
	}
}

// IsIsolated returns true when the port has at least medium isolation.
func IsIsolated(p scanner.PortInfo, peers []scanner.PortInfo) bool {
	return IsolationFor(p, peers) >= IsolationMedium
}
