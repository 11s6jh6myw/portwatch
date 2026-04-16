package eventlog

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Event represents a single port change event.
type Event struct {
	Timestamp time.Time      `json:"timestamp"`
	Kind      string         `json:"kind"` // "opened" | "closed"
	Port      scanner.PortInfo `json:"port"`
}

// Store appends events to a newline-delimited JSON log file.
type Store struct {
	mu   sync.Mutex
	path string
}

// NewStore returns a Store that writes to path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Append writes one event to the log.
func (s *Store) Append(kind string, port scanner.PortInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	e := Event{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Port:      port,
	}

	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(e)
}

// ReadAll returns all events from the log.
func (s *Store) ReadAll() ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var events []Event
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Event
		if err := dec.Decode(&e); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}
