package portmeta

import (
	"math"

	"github.com/user/portwatch/internal/scanner"
)

// ProximityLevel describes how close a port is to other known significant ports.
type ProximityLevel int

const (
	ProximityNone    ProximityLevel = iota // no nearby significant ports
	ProximityDistant                       // significant port >100 away
	ProximityNear                          // significant port within 100
	ProximityAdjacent                      // significant port within 10
	ProximityImmediate                     // significant port within 2
)

func (p ProximityLevel) String() string {
	switch p {
	case ProximityImmediate:
		return "immediate"
	case ProximityAdjacent:
		return "adjacent"
	case ProximityNear:
		return "near"
	case ProximityDistant:
		return "distant"
	default:
		return "none"
	}
}

// ProximityFor returns the proximity level of port p relative to the set of
// all open ports. It finds the minimum distance to any other significant
// (well-known) port in the provided slice.
func ProximityFor(p scanner.PortInfo, all []scanner.PortInfo) ProximityLevel {
	minDist := math.MaxInt32
	for _, other := range all {
		if other.Port == p.Port {
			continue
		}
		if _, ok := wellKnownPorts[other.Port]; !ok {
			continue
		}
		d := p.Port - other.Port
		if d < 0 {
			d = -d
		}
		if d < minDist {
			minDist = d
		}
	}
	if minDist == math.MaxInt32 {
		return ProximityNone
	}
	switch {
	case minDist <= 2:
		return ProximityImmediate
	case minDist <= 10:
		return ProximityAdjacent
	case minDist <= 100:
		return ProximityNear
	default:
		return ProximityDistant
	}
}
