// Package state manages on-disk persistence of port scan snapshots for
// portwatch. It serialises and deserialises []scanner.PortInfo values as
// JSON, enabling the daemon to compare the current port landscape against
// the last known good state even after a process restart.
//
// Typical usage:
//
//	store := state.NewStore("/var/lib/portwatch/state.json")
//
//	// Persist the current scan.
//	if err := store.Save(ports); err != nil { ... }
//
//	// Recover state on startup.
//	snap, err := store.Load()
//	if err != nil { ... }
//	prev := snap.Ports
package state
