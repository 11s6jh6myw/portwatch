package portmeta

// RiskLevel represents the severity of a port risk assessment.
type RiskLevel int

const (
	RiskNone RiskLevel = iota
	RiskLow
	RiskMedium
	RiskHigh
)

// String returns a human-readable label for the risk level.
func (r RiskLevel) String() string {
	switch r {
	case RiskLow:
		return "low"
	case RiskMedium:
		return "medium"
	case RiskHigh:
		return "high"
	default:
		return "none"
	}
}

// Score returns a RiskLevel for the given port number.
// Ports with known dangerous services score higher.
func Score(port int) RiskLevel {
	switch {
	case isHighRisk(port):
		return RiskHigh
	case isMediumRisk(port):
		return RiskMedium
	case IsRisky(port):
		return RiskLow
	default:
		return RiskNone
	}
}

func isHighRisk(port int) bool {
	high := map[int]bool{
		23:   true, // telnet
		512:  true, // rexec
		513:  true, // rlogin
		514:  true, // rsh
		1433: true, // mssql
		3389: true, // rdp
		5900: true, // vnc
	}
	return high[port]
}

func isMediumRisk(port int) bool {
	medium := map[int]bool{
		21:   true, // ftp
		69:   true, // tftp
		161:  true, // snmp
		3306: true, // mysql
		5432: true, // postgres
		6379: true, // redis
		9200: true, // elasticsearch
	}
	return medium[port]
}
