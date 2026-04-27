package portmeta

import "github.com/iamcalledrob/portwatch/internal/scanner"

// AffinityLevel describes how closely a port is associated with a known
// service family (e.g. web, database, messaging).
type AffinityLevel int

const (
	AffinityNone   AffinityLevel = iota // no recognised family
	AffinityWeak                        // loosely associated
	AffinityMedium                      // commonly associated
	AffinityStrong                      // canonical member of the family
)

func (a AffinityLevel) String() string {
	switch a {
	case AffinityWeak:
		return "weak"
	case AffinityMedium:
		return "medium"
	case AffinityStrong:
		return "strong"
	default:
		return "none"
	}
}

// serviceFamily maps a port number to a broad service family name.
var serviceFamily = map[int]string{
	80: "web", 443: "web", 8080: "web", 8443: "web",
	3306: "database", 5432: "database", 1433: "database", 27017: "database",
	5672: "messaging", 5671: "messaging", 4369: "messaging",
	6379: "cache", 11211: "cache",
	22: "remote", 23: "remote", 3389: "remote",
	25: "mail", 587: "mail", 465: "mail", 143: "mail", 993: "mail",
	53: "dns", 853: "dns",
}

// strongAffinity lists ports that are canonical members of their family.
var strongAffinity = map[int]bool{
	80: true, 443: true, 3306: true, 5432: true,
	22: true, 25: true, 53: true, 6379: true,
}

// AffinityFor returns the affinity level and family name for p.
// family is empty when AffinityNone is returned.
func AffinityFor(p scanner.PortInfo) (AffinityLevel, string) {
	family, ok := serviceFamily[p.Port]
	if !ok {
		return AffinityNone, ""
	}
	if strongAffinity[p.Port] {
		return AffinityStrong, family
	}
	// Ports sharing a family but not canonical get Medium; alternates get Weak.
	if p.Port < 1024 {
		return AffinityMedium, family
	}
	return AffinityWeak, family
}
