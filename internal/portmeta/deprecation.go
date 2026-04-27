package portmeta

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// DeprecationLevel indicates how deprecated or legacy a port's usage is considered.
type DeprecationLevel int

const (
	DeprecationNone DeprecationLevel = iota
	DeprecationMinor
	DeprecationModerate
	DeprecationHigh
	DeprecationObsolete
)

func (d DeprecationLevel) String() string {
	switch d {
	case DeprecationNone:
		return "none"
	case DeprecationMinor:
		return "minor"
	case DeprecationModerate:
		return "moderate"
	case DeprecationHigh:
		return "high"
	case DeprecationObsolete:
		return "obsolete"
	default:
		return "unknown"
	}
}

// obsoletePorts maps port numbers to the year they were deprecated or considered legacy.
var obsoletePorts = map[int]int{
	21:   1990, // FTP — largely superseded by SFTP/SCP
	23:   1995, // Telnet — replaced by SSH
	69:   2000, // TFTP — considered insecure
	79:   1998, // Finger — privacy concerns
	110:  2005, // POP3 — largely replaced by IMAP/HTTPS
	119:  2000, // NNTP — Usenet decline
	135:  2003, // MS-RPC — attack surface concerns
	139:  2003, // NetBIOS — replaced by SMB over TCP
	512:  2000, // rexec — insecure remote exec
	513:  2000, // rlogin — replaced by SSH
	514:  2000, // rsh — replaced by SSH
	2049: 2010, // NFSv2/v3 — NFSv4 preferred
}

// DeprecationFor returns the deprecation level for the given port.
func DeprecationFor(p scanner.PortInfo) DeprecationLevel {
	year, ok := obsoletePorts[p.Port]
	if !ok {
		return DeprecationNone
	}
	currentYear := time.Now().Year()
	age := currentYear - year
	switch {
	case age >= 30:
		return DeprecationObsolete
	case age >= 20:
		return DeprecationHigh
	case age >= 10:
		return DeprecationModerate
	case age >= 5:
		return DeprecationMinor
	default:
		return DeprecationNone
	}
}

// IsDeprecated returns true if the port has any non-zero deprecation level.
func IsDeprecated(p scanner.PortInfo) bool {
	return DeprecationFor(p) > DeprecationNone
}
