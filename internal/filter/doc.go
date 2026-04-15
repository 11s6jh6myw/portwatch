// Package filter implements port exclusion rules for portwatch.
//
// A Filter is constructed from a list of rule strings, where each rule
// is either a single port number (e.g. "22") or an inclusive port range
// (e.g. "8000-9000"). Once built, a Filter can be queried to determine
// whether a given port should be excluded from monitoring.
//
// Example usage:
//
//	f, err := filter.New([]string{"22", "80", "8000-9000"})
//	if err != nil {
//		log.Fatal(err)
//	}
//	if f.Excluded(8080) {
//		fmt.Println("port 8080 is filtered out")
//	}
package filter
