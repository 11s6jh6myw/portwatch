package portmeta

// ExposureLevel describes how exposed a port is considered to be.
type ExposureLevel int

const (
	ExposureNone    ExposureLevel = iota
	ExposureLow                   // internal / loopback only
	ExposureMedium                // LAN-facing
	ExposureHigh                  // internet-facing or widely scanned
)

func (e ExposureLevel) String() string {
	switch e {
	case ExposureLow:
		return "low"
	case ExposureMedium:
		return "medium"
	case ExposureHigh:
		return "high"
	default:
		return "none"
	}
}

// highExposurePorts are ports frequently targeted by internet scanners.
var highExposurePorts = map[int]bool{
	21: true, 22: true, 23: true, 25: true, 80: true,
	443: true, 3389: true, 8080: true, 8443: true, 3306: true,
	5432: true, 6379: true, 27017: true,
}

// mediumExposurePorts are ports common on internal networks.
var mediumExposurePorts = map[int]bool{
	139: true, 445: true, 111: true, 2049: true, 5900: true,
	5985: true, 5986: true, 9200: true, 9300: true,
}

// ExposureFor returns the ExposureLevel for a given port number.
func ExposureFor(port int) ExposureLevel {
	if highExposurePorts[port] {
		return ExposureHigh
	}
	if mediumExposurePorts[port] {
		return ExposureMedium
	}
	if port > 0 && port <= 1024 {
		return ExposureLow
	}
	return ExposureNone
}
