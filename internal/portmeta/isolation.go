package portmeta

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// IsolationLevel describes how isolated a port is from other active ports.
type IsolationLevel int

const (
	IsolationNone    IsolationLevel = iota // port has many neighbours
	IsolationLow                           // a few nearby ports are open
	IsolationMedium                        // sparse neighbourhood
	IsolationHigh                          // stands alone in its range
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

// IsolationFor returns how isolated port p is given the full set of currently
// open ports and the time the port was first seen.
func IsolationFor(p scanner.PortInfo, peers []scanner.PortInfo, firstSeen time.Time) IsolationLevel {
	neighbours := 0
	for _, peer := range peers {
		if peer.Port == p.Port {
			continue
		}
		diff := int(p.Port) - int(peer.Port)
		if diff < 0 {
			diff = -diff
		}
		if diff <= 100 {
			neighbours++
		}
	}

	switch {
	case neighbours == 0:
		return IsolationHigh
	case neighbours <= 2:
		return IsolationMedium
	case neighbours <= 5:
		return IsolationLow
	default:
		return IsolationNone
	}
}
