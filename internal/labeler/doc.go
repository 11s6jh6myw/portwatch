// Package labeler maps port numbers to human-readable service labels.
//
// It ships with a built-in table of well-known port-to-service mappings
// (ssh, http, postgres, etc.) and accepts caller-supplied overrides that
// take precedence over the defaults.
//
// The Enricher helper decorates scanner.PortInfo slices in place,
// attaching both a Label string and a Known boolean so downstream
// consumers can distinguish recognised ports from unexpected ones.
package labeler
