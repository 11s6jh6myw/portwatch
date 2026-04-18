package portmeta

// Category groups a port into a broad service class.
type Category string

const (
	CategoryDatabase  Category = "database"
	CategoryWeb       Category = "web"
	CategoryRemote    Category = "remote-access"
	CategoryMessaging Category = "messaging"
	CategoryDNS       Category = "dns"
	CategoryMail      Category = "mail"
	CategoryUnknown   Category = "unknown"
)

// String returns the string representation of a Category.
func (c Category) String() string { return string(c) }

var portCategories = map[uint16]Category{
	21:   CategoryRemote,
	22:   CategoryRemote,
	23:   CategoryRemote,
	25:   CategoryMail,
	53:   CategoryDNS,
	80:   CategoryWeb,
	110:  CategoryMail,
	143:  CategoryMail,
	443:  CategoryWeb,
	3306: CategoryDatabase,
	5432: CategoryDatabase,
	6379: CategoryDatabase,
	27017: CategoryDatabase,
	5672:  CategoryMessaging,
	9092:  CategoryMessaging,
	8080:  CategoryWeb,
	8443:  CategoryWeb,
	3389:  CategoryRemote,
	5900:  CategoryRemote,
}

// Categorize returns the Category for the given port number.
// Returns CategoryUnknown if the port is not recognised.
func Categorize(port uint16) Category {
	if c, ok := portCategories[port]; ok {
		return c
	}
	return CategoryUnknown
}

// CategorizeAll returns a map of port → Category for the supplied ports.
func CategorizeAll(ports []uint16) map[uint16]Category {
	out := make(map[uint16]Category, len(ports))
	for _, p := range ports {
		out[p] = Categorize(p)
	}
	return out
}
