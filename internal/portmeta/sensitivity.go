package portmeta

// SensitivityLevel indicates how sensitive a port's traffic is considered.
type SensitivityLevel int

const (
	SensitivityNone SensitivityLevel = iota
	SensitivityLow
	SensitivityMedium
	SensitivityHigh
	SensitivityCritical
)

func (s SensitivityLevel) String() string {
	switch s {
	case SensitivityLow:
		return "low"
	case SensitivityMedium:
		return "medium"
	case SensitivityHigh:
		return "high"
	case SensitivityCritical:
		return "critical"
	default:
		return "none"
	}
}

// sensitivityMap maps well-known ports to their sensitivity level.
var sensitivityMap = map[uint16]SensitivityLevel{
	21:   SensitivityHigh,     // FTP
	22:   SensitivityHigh,     // SSH
	23:   SensitivityCritical, // Telnet
	25:   SensitivityMedium,   // SMTP
	53:   SensitivityMedium,   // DNS
	80:   SensitivityLow,      // HTTP
	110:  SensitivityMedium,   // POP3
	143:  SensitivityMedium,   // IMAP
	389:  SensitivityHigh,     // LDAP
	443:  SensitivityLow,      // HTTPS
	445:  SensitivityCritical, // SMB
	1433: SensitivityCritical, // MSSQL
	3306: SensitivityCritical, // MySQL
	3389: SensitivityCritical, // RDP
	5432: SensitivityCritical, // PostgreSQL
	5900: SensitivityHigh,     // VNC
	6379: SensitivityHigh,     // Redis
	27017: SensitivityHigh,    // MongoDB
}

// SensitivityFor returns the sensitivity level for the given port.
func SensitivityFor(port uint16) SensitivityLevel {
	if s, ok := sensitivityMap[port]; ok {
		return s
	}
	return SensitivityNone
}

// IsSensitive returns true if the port has at least medium sensitivity.
func IsSensitive(port uint16) bool {
	return SensitivityFor(port) >= SensitivityMedium
}
