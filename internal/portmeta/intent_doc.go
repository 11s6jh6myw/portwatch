// Package portmeta provides metadata enrichment for scanned ports.
//
// The intent sub-feature infers the operational intent of an open port —
// for example, whether it serves infrastructure, application, development,
// administrative, or legacy purposes.
//
// Usage:
//
//	annotate := portmeta.NewIntentAnnotator()
//	enriched := annotate(ports)
//
// Ports may also be filtered by intent level using FilterByIntent.
package portmeta
