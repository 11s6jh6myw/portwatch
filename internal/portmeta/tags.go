package portmeta

// Tag represents a descriptive label attached to a port.
type Tag string

const (
	TagDatabase  Tag = "database"
	TagWeb       Tag = "web"
	TagRemote    Tag = "remote-access"
	TagMail      Tag = "mail"
	TagDNS       Tag = "dns"
	TagFile      Tag = "file-transfer"
	TagMonitor   Tag = "monitoring"
	TagContainer Tag = "container"
	TagUnknown   Tag = "unknown"
)

// portTags maps well-known ports to descriptive tags.
var portTags = map[uint16][]Tag{
	21:    {TagFile},
	22:    {TagRemote},
	23:    {TagRemote},
	25:    {TagMail},
	53:    {TagDNS},
	80:    {TagWeb},
	443:   {TagWeb},
	3306:  {TagDatabase},
	5432:  {TagDatabase},
	6379:  {TagDatabase},
	8080:  {TagWeb},
	8443:  {TagWeb},
	27017: {TagDatabase},
	2375:  {TagContainer},
	2376:  {TagContainer},
	9090:  {TagMonitor},
	9100:  {TagMonitor},
}

// TagsFor returns the tags associated with the given port number.
// If no tags are defined, []Tag{TagUnknown} is returned.
func TagsFor(port uint16) []Tag {
	if tags, ok := portTags[port]; ok {
		return tags
	}
	return []Tag{TagUnknown}
}

// HasTag reports whether the given port carries the specified tag.
func HasTag(port uint16, tag Tag) bool {
	for _, t := range TagsFor(port) {
		if t == tag {
			return true
		}
	}
	return false
}
