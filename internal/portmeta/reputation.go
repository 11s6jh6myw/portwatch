package portmeta

// ReputationLevel indicates how well-regarded a port is in the security community.
type ReputationLevel int

const (
	ReputationUnknown ReputationLevel = iota
	ReputationPoor
	ReputationFair
	ReputationGood
	ReputationTrusted
)

func (r ReputationLevel) String() string {
	switch r {
	case ReputationPoor:
		return "poor"
	case ReputationFair:
		return "fair"
	case ReputationGood:
		return "good"
	case ReputationTrusted:
		return "trusted"
	default:
		return "unknown"
	}
}

// reputationMap maps well-known ports to their reputation level.
var reputationMap = map[uint16]ReputationLevel{
	22:   ReputationGood,
	80:   ReputationTrusted,
	443:  ReputationTrusted,
	8080: ReputationFair,
	8443: ReputationGood,
	23:   ReputationPoor,
	21:   ReputationFair,
	25:   ReputationFair,
	3306: ReputationGood,
	5432: ReputationGood,
	6379: ReputationFair,
	27017: ReputationFair,
	4444:  ReputationPoor,
	1337:  ReputationPoor,
	31337: ReputationPoor,
}

// ReputationFor returns the reputation level for the given port number.
func ReputationFor(port uint16) ReputationLevel {
	if r, ok := reputationMap[port]; ok {
		return r
	}
	return ReputationUnknown
}

// IsReputable returns true if the port has a Fair reputation or better.
func IsReputable(port uint16) bool {
	return ReputationFor(port) >= ReputationFair
}
