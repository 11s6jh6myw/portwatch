package portmeta

// CriticalityLevel represents how critical a port is to system operation.
type CriticalityLevel int

const (
	CriticalityNone CriticalityLevel = iota
	CriticalityLow
	CriticalityMedium
	CriticalityHigh
	CriticalityCritical
)

func (c CriticalityLevel) String() string {
	switch c {
	case CriticalityLow:
		return "low"
	case CriticalityMedium:
		return "medium"
	case CriticalityHigh:
		return "high"
	case CriticalityCritical:
		return "critical"
	default:
		return "none"
	}
}

// criticalityMap maps well-known ports to their criticality level.
var criticalityMap = map[uint16]CriticalityLevel{
	22:   CriticalityHigh,     // SSH
	25:   CriticalityHigh,     // SMTP
	53:   CriticalityCritical, // DNS
	80:   CriticalityMedium,   // HTTP
	443:  CriticalityHigh,     // HTTPS
	3306: CriticalityHigh,     // MySQL
	5432: CriticalityHigh,     // PostgreSQL
	6379: CriticalityHigh,     // Redis
	27017: CriticalityHigh,    // MongoDB
	2379: CriticalityCritical, // etcd
	6443: CriticalityCritical, // Kubernetes API
	8080: CriticalityMedium,   // HTTP alt
	8443: CriticalityMedium,   // HTTPS alt
	23:   CriticalityLow,      // Telnet
	21:   CriticalityMedium,   // FTP
}

// CriticalityFor returns the criticality level for the given port number.
func CriticalityFor(port uint16) CriticalityLevel {
	if c, ok := criticalityMap[port]; ok {
		return c
	}
	return CriticalityNone
}

// IsCritical returns true if the port is high or critical.
func IsCritical(port uint16) bool {
	c := CriticalityFor(port)
	return c >= CriticalityHigh
}
