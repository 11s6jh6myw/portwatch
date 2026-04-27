// Package portmeta — signal.go
//
// SignalStrength provides a single composite indicator that combines
// risk, anomaly, urgency, and recency into one actionable value.
//
// Levels (ascending severity):
//
//	none     – no meaningful signal detected
//	weak     – minor indicator; low priority
//	moderate – worth investigating
//	strong   – high-confidence indicator; prompt review advised
//	critical – immediate action recommended
//
// Use NewSignalAnnotator to attach signal metadata to a port slice, and
// FilterByMinSignal to narrow results to a minimum strength threshold.
package portmeta
