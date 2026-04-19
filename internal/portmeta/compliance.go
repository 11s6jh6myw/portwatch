package portmeta

// ComplianceFlag indicates whether a port is relevant to common compliance frameworks.
type ComplianceFlag string

const (
	CompliancePCI  ComplianceFlag = "pci"
	ComplianceHIPAA ComplianceFlag = "hipaa"
	ComplianceSOC2 ComplianceFlag = "soc2"
	ComplianceNone ComplianceFlag = "none"
)

var complianceMap = map[int][]ComplianceFlag{
	22:    {CompliancePCI, ComplianceSOC2},
	23:    {CompliancePCI, ComplianceHIPAA, ComplianceSOC2},
	80:    {CompliancePCI},
	443:   {CompliancePCI, ComplianceSOC2},
	3306:  {CompliancePCI, ComplianceHIPAA},
	5432:  {CompliancePCI, ComplianceHIPAA},
	6379:  {CompliancePCI},
	27017: {CompliancePCI, ComplianceHIPAA},
}

// ComplianceFlagsFor returns all compliance flags associated with a port.
func ComplianceFlagsFor(port int) []ComplianceFlag {
	if flags, ok := complianceMap[port]; ok {
		return flags
	}
	return []ComplianceFlag{ComplianceNone}
}

// HasComplianceFlag returns true if the port is tagged with the given flag.
func HasComplianceFlag(port int, flag ComplianceFlag) bool {
	for _, f := range ComplianceFlagsFor(port) {
		if f == flag {
			return true
		}
	}
	return false
}
