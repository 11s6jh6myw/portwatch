package portmeta

// Dependency describes a known relationship between two ports.
type Dependency struct {
	Port    int
	Related int
	Reason  string
}

// knownDependencies maps a port to ports it commonly depends on or is paired with.
var knownDependencies = []Dependency{
	{Port: 443, Related: 80, Reason: "HTTPS often paired with HTTP"},
	{Port: 8443, Related: 8080, Reason: "Alt HTTPS paired with alt HTTP"},
	{Port: 3306, Related: 33060, Reason: "MySQL paired with MySQL X Protocol"},
	{Port: 5432, Related: 5433, Reason: "PostgreSQL primary and replica"},
	{Port: 6379, Related: 26379, Reason: "Redis paired with Sentinel"},
	{Port: 9200, Related: 9300, Reason: "Elasticsearch HTTP and transport"},
	{Port: 2181, Related: 2888, Reason: "ZooKeeper client and peer"},
	{Port: 2888, Related: 3888, Reason: "ZooKeeper peer and leader election"},
}

// DependenciesFor returns all known dependencies for the given port.
func DependenciesFor(port int) []Dependency {
	var out []Dependency
	for _, d := range knownDependencies {
		if d.Port == port || d.Related == port {
			out = append(out, d)
		}
	}
	return out
}

// HasDependency reports whether port and related are known to be paired.
func HasDependency(port, related int) bool {
	for _, d := range knownDependencies {
		if (d.Port == port && d.Related == related) ||
			(d.Port == related && d.Related == port) {
			return true
		}
	}
	return false
}
