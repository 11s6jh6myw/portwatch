// Package portmeta provides metadata enrichment and classification for
// observed network ports.
//
// # Lifespan
//
// The lifespan module classifies how long a port has been continuously
// observed open, based on its first-seen timestamp.
//
// Levels (ascending duration):
//
//	LifespanEphemeral  – open for less than one minute
//	LifespanShort      – open between one minute and one hour
//	LifespanMedium     – open between one hour and 24 hours
//	LifespanLong       – open between one day and one week
//	LifespanPermanent  – open for more than one week
//
// NewLifespanAnnotator returns a port-slice transformer that reads the
// "first_seen" meta field and writes the resolved lifespan level string
// into the "lifespan" meta field.
//
// FilterByMinLifespan discards ports whose lifespan falls below the
// specified minimum. Ports with no lifespan metadata are always kept.
package portmeta
