package portmeta

// VisibilityLevel describes how publicly known or accessible a port is.
type VisibilityLevel int

const (
	VisibilityUnknown VisibilityLevel = iota
	VisibilityInternal
	VisibilityRestricted
	VisibilityPublic
)

func (v VisibilityLevel) String() string {
	switch v {
	case VisibilityInternal:
		return "internal"
	case VisibilityRestricted:
		return "restricted"
	case VisibilityPublic:
		return "public"
	default:
		return "unknown"
	}
}

// VisibilityFor returns the expected visibility level for a given port.
func VisibilityFor(port int) VisibilityLevel {
	switch {
	case isPublicPort(port):
		return VisibilityPublic
	case isRestrictedPort(port):
		return VisibilityRestricted
	case isInternalPort(port):
		return VisibilityInternal
	default:
		return VisibilityUnknown
	}
}

func isPublicPort(port int) bool {
	public := []int{80, 443, 8080, 8443, 3000, 5000}
	for _, p := range public {
		if port == p {
			return true
		}
	}
	return false
}

func isRestrictedPort(port int) bool {
	restricted := []int{22, 3306, 5432, 6379, 27017, 9200}
	for _, p := range restricted {
		if port == p {
			return true
		}
	}
	return false
}

func isInternalPort(port int) bool {
	return port >= 49152 && port <= 65535
}
