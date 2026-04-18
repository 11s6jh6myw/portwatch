package portmeta

// Owner holds attribution metadata for a well-known port.
type Owner struct {
	Org     string // organisation or project responsible
	Contact string // URL or email for more information
}

// ownerMap maps port numbers to known owners.
var ownerMap = map[uint16]Owner{
	22:   {Org: "IETF", Contact: "https://tools.ietf.org/html/rfc4251"},
	25:   {Org: "IETF", Contact: "https://tools.ietf.org/html/rfc5321"},
	53:   {Org: "IETF", Contact: "https://tools.ietf.org/html/rfc1035"},
	80:   {Org: "IETF", Contact: "https://tools.ietf.org/html/rfc7230"},
	443:  {Org: "IETF", Contact: "https://tools.ietf.org/html/rfc8446"},
	3306: {Org: "Oracle", Contact: "https://www.mysql.com"},
	5432: {Org: "PostgreSQL Global Development Group", Contact: "https://www.postgresql.org"},
	6379: {Org: "Redis Ltd", Contact: "https://redis.io"},
	27017: {Org: "MongoDB Inc", Contact: "https://www.mongodb.com"},
	9200: {Org: "Elastic NV", Contact: "https://www.elastic.co"},
}

// LookupOwner returns the Owner for a given port and whether one was found.
func LookupOwner(port uint16) (Owner, bool) {
	o, ok := ownerMap[port]
	return o, ok
}

// KnownOwner reports whether a port has a registered owner.
func KnownOwner(port uint16) bool {
	_, ok := ownerMap[port]
	return ok
}
