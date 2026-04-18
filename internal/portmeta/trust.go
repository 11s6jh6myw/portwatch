package portmeta

// TrustLevel indicates how trusted a port/service is considered.
type TrustLevel int

const (
	TrustUnknown TrustLevel = iota
	TrustLow
	TrustMedium
	TrustHigh
)

func (t TrustLevel) String() string {
	switch t {
	case TrustLow:
		return "low"
	case TrustMedium:
		return "medium"
	case TrustHigh:
		return "high"
	default:
		return "unknown"
	}
}

// trustTable maps well-known ports to a trust level.
var trustTable = map[uint16]TrustLevel{
	22:   TrustHigh,   // SSH
	80:   TrustHigh,   // HTTP
	443:  TrustHigh,   // HTTPS
	53:   TrustMedium, // DNS
	25:   TrustMedium, // SMTP
	3306: TrustMedium, // MySQL
	5432: TrustMedium, // PostgreSQL
	23:   TrustLow,    // Telnet
	21:   TrustLow,    // FTP
	4444: TrustLow,    // common backdoor
	1337: TrustLow,
}

// TrustFor returns the TrustLevel for the given port number.
func TrustFor(port uint16) TrustLevel {
	if t, ok := trustTable[port]; ok {
		return t
	}
	return TrustUnknown
}

// IsTrusted returns true if the port has at least TrustMedium trust.
func IsTrusted(port uint16) bool {
	t := TrustFor(port)
	return t == TrustMedium || t == TrustHigh
}
