// Package audit provides a structured audit trail for port change events,
// recording who observed a change, when, and what action was taken.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Hostname  string    `json:"hostname"`
	Kind      string    `json:"kind"` // "opened" | "closed"
	Port      uint16    `json:"port"`
	Protocol  string    `json:"protocol"`
	Action    string    `json:"action"` // "alerted" | "suppressed" | "filtered"
}

// Store appends audit entries to a newline-delimited JSON file.
type Store struct {
	mu   sync.Mutex
	path string
}

// NewStore returns a Store that writes to the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Record appends an Entry to the audit log.
func (s *Store) Record(e Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open file: %w", err)
	}
	defer f.Close()

	line, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(f, "%s\n", line)
	return err
}

// ReadAll returns all entries from the audit log.
func (s *Store) ReadAll() ([]Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: read file: %w", err)
	}

	var entries []Entry
	dec := json.NewDecoder(
		// reuse bytes as a reader
		newBytesReader(data),
	)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
