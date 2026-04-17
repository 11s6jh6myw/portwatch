// Package tagset provides a lightweight key=value label system for port events.
//
// A [Set] holds an ordered, deduplicated collection of string tags that can be
// attached to any event or snapshot to carry environment metadata (host, region,
// team, etc.).
//
// Tags can be supplied explicitly via [New] or enriched automatically from the
// process environment using [Resolver]. Environment variables prefixed with
// PORTWATCH_TAG_ are lowercased and injected as tags; the host name is always
// added when available.
//
// Example:
//
//	s := tagset.New("env=prod", "team=infra")
//	r := tagset.NewResolver()
//	enriched := r.Enrich(s)
//	fmt.Println(enriched) // a=1,env=prod,host=myhost,team=infra
package tagset
