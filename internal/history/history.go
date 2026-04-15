// Package history records port change events over time,
// allowing portwatch to provide a persistent audit trail
// of opened and closed ports.
package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// EventType indicates whether a port was opened or closed.
type EventType string

const (
	EventOpened EventType = "opened"
	EventClosed EventType = "closed"
)

// Event represents a single port change occurrence.
type Event struct {
	Timestamp time.Time      `json:"timestamp"`
	Type      EventType      `json:"type"`
	Port      scanner.PortInfo `json:"port"`
}

// Store persists and retrieves port change history.
type Store struct {
	mu     sync.Mutex
	path   string
	events []Event
}

// NewStore creates a Store backed by the given file path.
// Existing events are loaded if the file is present.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// Record appends a new event and flushes it to disk.
func (s *Store) Record(eventType EventType, port scanner.PortInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = append(s.events, Event{
		Timestamp: time.Now().UTC(),
		Type:      eventType,
		Port:      port,
	})
	return s.save()
}

// Events returns a copy of all recorded events.
func (s *Store) Events() []Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := make([]Event, len(s.events))
	for i, e := range s.events {
		copy[i] = e
	}
	return copy
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.events)
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.events, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
