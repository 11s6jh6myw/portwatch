// Package portmeta provides a built-in lookup table of well-known TCP port
// numbers with associated service names, protocol hints, and risk flags.
//
// It is used by the labeler and reporter packages to enrich port scan results
// with human-readable context and to surface potentially dangerous open ports.
//
// Example:
//
//	if m, ok := portmeta.Lookup(3306); ok {
//		fmt.Println(m) // "3306/tcp (MySQL)"
//	}
package portmeta
