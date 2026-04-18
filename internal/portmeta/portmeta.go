// Package portmeta provides metadata lookups for well-known port numbers,
// including protocol hints and common service descriptions.
package portmeta

import "fmt"

// Meta holds descriptive metadata for a port number.
type Meta struct {
	Port        uint16
	Protocol    string // e.g. "tcp"
	Service     string // e.g. "HTTP"
	Description string
	Risky       bool // true if commonly exploited
}

// String returns a short human-readable summary.
func (m Meta) String() string {
	return fmt.Sprintf("%d/%s (%s)", m.Port, m.Protocol, m.Service)
}

var table = map[uint16]Meta{
	21:   {21, "tcp", "FTP", "File Transfer Protocol", true},
	22:   {22, "tcp", "SSH", "Secure Shell", false},
	23:   {23, "tcp", "Telnet", "Unencrypted remote shell", true},
	25:   {25, "tcp", "SMTP", "Simple Mail Transfer Protocol", false},
	53:   {53, "tcp", "DNS", "Domain Name System", false},
	80:   {80, "tcp", "HTTP", "Hypertext Transfer Protocol", false},
	110:  {110, "tcp", "POP3", "Post Office Protocol v3", false},
	143:  {143, "tcp", "IMAP", "Internet Message Access Protocol", false},
	443:  {443, "tcp", "HTTPS", "HTTP over TLS", false},
	445:  {445, "tcp", "SMB", "Server Message Block", true},
	3306: {3306, "tcp", "MySQL", "MySQL Database", true},
	3389: {3389, "tcp", "RDP", "Remote Desktop Protocol", true},
	5432: {5432, "tcp", "PostgreSQL", "PostgreSQL Database", false},
	6379: {6379, "tcp", "Redis", "Redis in-memory store", true},
	8080: {8080, "tcp", "HTTP-Alt", "Alternate HTTP", false},
	8443: {8443, "tcp", "HTTPS-Alt", "Alternate HTTPS", false},
	27017: {27017, "tcp", "MongoDB", "MongoDB Database", true},
}

// Lookup returns metadata for the given port number.
// The second return value is false when the port is not in the built-in table.
func Lookup(port uint16) (Meta, bool) {
	m, ok := table[port]
	return m, ok
}

// IsRisky reports whether the port is flagged as commonly exploited.
func IsRisky(port uint16) bool {
	m, ok := table[port]
	return ok && m.Risky
}
