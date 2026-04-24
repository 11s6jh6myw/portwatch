package portmeta

// ReachabilityLevel describes how accessible a port is from external networks.
type ReachabilityLevel int

const (
	ReachabilityNone    ReachabilityLevel = iota // not reachable externally
	ReachabilityPrivate                          // reachable within private networks
	ReachabilityLimited                          // reachable with authentication or firewall rules
	ReachabilityPublic                           // openly reachable from the internet
)

func (r ReachabilityLevel) String() string {
	switch r {
	case ReachabilityNone:
		return "none"
	case ReachabilityPrivate:
		return "private"
	case ReachabilityLimited:
		return "limited"
	case ReachabilityPublic:
		return "public"
	default:
		return "unknown"
	}
}

// publicPorts are ports typically exposed to the internet.
var publicPorts = map[int]struct{}{
	80: {}, 443: {}, 8080: {}, 8443: {}, 21: {}, 22: {}, 25: {}, 53: {},
	110: {}, 143: {}, 3000: {}, 5000: {}, 8000: {},
}

// privatePorts are ports commonly used in internal infrastructure.
var privatePorts = map[int]struct{}{
	5432: {}, 3306: {}, 27017: {}, 6379: {}, 9200: {}, 5601: {},
	2181: {}, 9092: {}, 8500: {}, 4369: {},
}

// limitedPorts are ports that are reachable but typically gated.
var limitedPorts = map[int]struct{}{
	3389: {}, 5900: {}, 2222: {}, 8888: {}, 9090: {}, 9443: {}, 15672: {},
}

// ReachabilityFor returns the estimated reachability level for a given port number.
func ReachabilityFor(port int) ReachabilityLevel {
	if _, ok := publicPorts[port]; ok {
		return ReachabilityPublic
	}
	if _, ok := limitedPorts[port]; ok {
		return ReachabilityLimited
	}
	if _, ok := privatePorts[port]; ok {
		return ReachabilityPrivate
	}
	return ReachabilityNone
}

// IsReachable returns true if the port has any reachability above None.
func IsReachable(port int) bool {
	return ReachabilityFor(port) > ReachabilityNone
}
