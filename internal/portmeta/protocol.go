package portmeta

// Protocol represents the transport protocol for a port.
type Protocol string

const (
	ProtocolTCP     Protocol = "tcp"
	ProtocolUDP     Protocol = "udp"
	ProtocolUnknown Protocol = "unknown"
)

// String returns the string representation of a Protocol.
func (p Protocol) String() string {
	return string(p)
}

// protocolMap maps well-known ports to their primary protocol.
var protocolMap = map[uint16]Protocol{
	20:   ProtocolTCP,
	21:   ProtocolTCP,
	22:   ProtocolTCP,
	23:   ProtocolTCP,
	25:   ProtocolTCP,
	53:   ProtocolUDP,
	67:   ProtocolUDP,
	68:   ProtocolUDP,
	80:   ProtocolTCP,
	110:  ProtocolTCP,
	123:  ProtocolUDP,
	143:  ProtocolTCP,
	161:  ProtocolUDP,
	443:  ProtocolTCP,
	465:  ProtocolTCP,
	514:  ProtocolUDP,
	587:  ProtocolTCP,
	993:  ProtocolTCP,
	995:  ProtocolTCP,
	3306: ProtocolTCP,
	3389: ProtocolTCP,
	5432: ProtocolTCP,
	5900: ProtocolTCP,
	6379: ProtocolTCP,
	8080: ProtocolTCP,
	8443: ProtocolTCP,
	27017: ProtocolTCP,
}

// ProtocolFor returns the primary protocol for a given port number.
// Returns ProtocolUnknown if the port is not recognised.
func ProtocolFor(port uint16) Protocol {
	if p, ok := protocolMap[port]; ok {
		return p
	}
	return ProtocolUnknown
}

// IsUDP reports whether the given port is primarily UDP.
func IsUDP(port uint16) bool {
	return ProtocolFor(port) == ProtocolUDP
}

// IsTCP reports whether the given port is primarily TCP.
func IsTCP(port uint16) bool {
	return ProtocolFor(port) == ProtocolTCP
}
