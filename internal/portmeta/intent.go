package portmeta

// IntentLevel describes the inferred purpose or intent of a port being open.
type IntentLevel int

const (
	IntentUnknown IntentLevel = iota
	IntentInfrastructure
	IntentApplication
	IntentDevelopment
	IntentAdministrative
	IntentLegacy
)

func (i IntentLevel) String() string {
	switch i {
	case IntentInfrastructure:
		return "infrastructure"
	case IntentApplication:
		return "application"
	case IntentDevelopment:
		return "development"
	case IntentAdministrative:
		return "administrative"
	case IntentLegacy:
		return "legacy"
	default:
		return "unknown"
	}
}

var intentMap = map[int]IntentLevel{
	22:   IntentAdministrative, // SSH
	23:   IntentLegacy,         // Telnet
	25:   IntentInfrastructure, // SMTP
	53:   IntentInfrastructure, // DNS
	80:   IntentApplication,
	443:  IntentApplication,
	3306: IntentApplication, // MySQL
	5432: IntentApplication, // PostgreSQL
	6379: IntentApplication, // Redis
	8080: IntentDevelopment,
	8443: IntentDevelopment,
	9090: IntentDevelopment,
	3000: IntentDevelopment,
	5000: IntentDevelopment,
	27017: IntentApplication, // MongoDB
}

// IntentFor returns the inferred intent level for a given port.
func IntentFor(port int) IntentLevel {
	if l, ok := intentMap[port]; ok {
		return l
	}
	return IntentUnknown
}
