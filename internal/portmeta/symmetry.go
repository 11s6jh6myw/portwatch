package portmeta

import (
	"math"

	"github.com/user/portwatch/internal/scanner"
)

// SymmetryLevel describes how symmetrically a port appears relative to its
// protocol peers (e.g. paired request/response ports, well-known siblings).
type SymmetryLevel int

const (
	SymmetryNone   SymmetryLevel = iota // no recognisable pairing
	SymmetryLow                         // weakly paired
	SymmetryMedium                      // moderately paired
	SymmetryHigh                        // strongly paired with a known sibling
)

func (s SymmetryLevel) String() string {
	switch s {
	case SymmetryLow:
		return "low"
	case SymmetryMedium:
		return "medium"
	case SymmetryHigh:
		return "high"
	default:
		return "none"
	}
}

// symmetryPairs maps a port to its canonical sibling port number.
var symmetryPairs = map[int]int{
	20: 21,   // FTP data / control
	21: 20,
	80: 443,  // HTTP / HTTPS
	443: 80,
	8080: 8443,
	8443: 8080,
	3306: 33060, // MySQL / MySQL X
	33060: 3306,
	5432: 5433, // PostgreSQL primary / replica
	5433: 5432,
	6379: 6380, // Redis primary / replica
	6380: 6379,
}

// SymmetryFor returns the symmetry level for port p given the full set of
// currently open ports. A port scores higher when its known sibling is also
// present in the open set.
func SymmetryFor(p scanner.PortInfo, open []scanner.PortInfo) SymmetryLevel {
	sibling, ok := symmetryPairs[p.Port]
	if !ok {
		return SymmetryNone
	}

	// Check whether the sibling is present.
	siblingOpen := false
	for _, o := range open {
		if o.Port == sibling {
			siblingOpen = true
			break
		}
	}

	if !siblingOpen {
		return SymmetryLow
	}

	// Score by port distance: closer siblings are more symmetric.
	dist := math.Abs(float64(p.Port - sibling))
	switch {
	case dist <= 1:
		return SymmetryHigh
	case dist <= 10:
		return SymmetryMedium
	default:
		return SymmetryLow
	}
}
