package portmeta

import "time"

// PrevalenceLevel indicates how commonly a port is observed across systems.
type PrevalenceLevel int

const (
	PrevalenceUnknown PrevalenceLevel = iota
	PrevalenceRare
	PrevalenceUncommon
	PrevalenceCommon
	PrevalenceUbiquitous
)

func (p PrevalenceLevel) String() string {
	switch p {
	case PrevalenceRare:
		return "rare"
	case PrevalenceUncommon:
		return "uncommon"
	case PrevalenceCommon:
		return "common"
	case PrevalenceUbiquitous:
		return "ubiquitous"
	default:
		return "unknown"
	}
}

// ubiquitousPorts are ports seen on virtually every networked system.
var ubiquitousPorts = map[int]bool{
	22: true, 80: true, 443: true,
}

// commonPorts are widely deployed but not universal.
var commonPorts = map[int]bool{
	21: true, 25: true, 53: true, 3306: true, 5432: true,
	6379: true, 8080: true, 8443: true,
}

// uncommonPorts are legitimate but niche services.
var uncommonPorts = map[int]bool{
	389: true, 636: true, 1433: true, 1521: true,
	5900: true, 9200: true, 27017: true,
}

// PrevalenceFor returns the prevalence level for a given port number.
// firstSeen is used to down-grade prevalence for very recently observed ports.
func PrevalenceFor(port int, firstSeen time.Time) PrevalenceLevel {
	if ubiquitousPorts[port] {
		return PrevalenceUbiquitous
	}
	if commonPorts[port] {
		return PrevalenceCommon
	}
	if uncommonPorts[port] {
		if !firstSeen.IsZero() && time.Since(firstSeen) < 24*time.Hour {
			return PrevalenceRare
		}
		return PrevalenceUncommon
	}
	return PrevalenceRare
}

// IsPrevalent returns true when the level is Common or higher.
func IsPrevalent(level PrevalenceLevel) bool {
	return level >= PrevalenceCommon
}
