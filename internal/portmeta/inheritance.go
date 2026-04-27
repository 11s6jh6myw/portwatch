package portmeta

import "github.com/user/portwatch/internal/scanner"

// InheritanceLevel describes how a port's behaviour is derived from a
// parent service or well-known ancestor.
type InheritanceLevel int

const (
	InheritanceNone InheritanceLevel = iota
	InheritanceDerived
	InheritanceExtended
	InheritanceCore
)

func (l InheritanceLevel) String() string {
	switch l {
	case InheritanceDerived:
		return "derived"
	case InheritanceExtended:
		return "extended"
	case InheritanceCore:
		return "core"
	default:
		return "none"
	}
}

// inheritanceMap maps port numbers to their inheritance level.
var inheritanceMap = map[uint16]InheritanceLevel{
	// Core internet infrastructure
	22:   InheritanceCore,
	25:   InheritanceCore,
	53:   InheritanceCore,
	80:   InheritanceCore,
	443:  InheritanceCore,
	// Extended from core protocols
	8080: InheritanceExtended,
	8443: InheritanceExtended,
	587:  InheritanceExtended,
	465:  InheritanceExtended,
	// Derived / application-layer
	3306: InheritanceDerived,
	5432: InheritanceDerived,
	6379: InheritanceDerived,
	27017: InheritanceDerived,
	9200: InheritanceDerived,
}

// InheritanceFor returns the InheritanceLevel for the given port.
func InheritanceFor(p scanner.PortInfo) InheritanceLevel {
	if lvl, ok := inheritanceMap[uint16(p.Port)]; ok {
		return lvl
	}
	return InheritanceNone
}

// IsInherited reports whether the port carries any inheritance classification.
func IsInherited(p scanner.PortInfo) bool {
	return InheritanceFor(p) != InheritanceNone
}
