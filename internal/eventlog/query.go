package eventlog

import "time"

// Filter holds optional constraints for querying events.
type Filter struct {
	Kind  string    // "opened", "closed", or "" for all
	Since time.Time // zero means no lower bound
	Port  int       // 0 means all ports
}

// Query returns events matching f from the store.
func (s *Store) Query(f Filter) ([]Event, error) {
	all, err := s.ReadAll()
	if err != nil {
		return nil, err
	}

	var out []Event
	for _, e := range all {
		if f.Kind != "" && e.Kind != f.Kind {
			continue
		}
		if !f.Since.IsZero() && e.Timestamp.Before(f.Since) {
			continue
		}
		if f.Port != 0 && e.Port.Port != f.Port {
			continue
		}
		out = append(out, e)
	}
	return out, nil
}
