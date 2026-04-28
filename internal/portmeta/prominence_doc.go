// Package portmeta — prominence
//
// ProminenceLevel captures how notable a port is within the observed
// network landscape. It combines signals from prevalence, criticality,
// impact, and recent activity to produce a single ordinal score.
//
// Levels (ascending):
//
//	none     – no supporting signals detected
//	low      – weakly supported by one minor signal
//	medium   – two or more minor signals present
//	high     – critical or high-impact port recently active
//	critical – multiple strong signals including recent activity
//
// Use NewProminenceAnnotator to attach the level to port metadata and
// FilterByMinProminence to narrow a port list to a minimum prominence.
package portmeta
