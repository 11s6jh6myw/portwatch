// Package history provides a persistent, append-only log of port change
// events detected by portwatch.
//
// Each time a port is observed to open or close, an Event is recorded
// with a timestamp and port details. Events are stored as JSON on disk
// and survive process restarts, giving operators a full audit trail of
// network activity over time.
//
// Usage:
//
//	store, err := history.NewStore("/var/lib/portwatch/history.json")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	store.Record(history.EventOpened, port)
//	for _, e := range store.Events() {
//	    fmt.Println(e.Timestamp, e.Type, e.Port)
//	}
package history
