// Package baseline manages the trusted port baseline used to detect
// unexpected changes during monitoring sessions.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot represents a saved baseline of known-good open ports.
type Snapshot struct {
	CreatedAt time.Time          `json:"created_at"`
	Ports     []scanner.PortInfo `json:"ports"`
}

// Store persists and retrieves baseline snapshots from disk.
type Store struct {
	path string
}

// NewStore creates a Store that reads/writes baselines at the given path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the provided ports as the current baseline snapshot.
func (s *Store) Save(ports []scanner.PortInfo) error {
	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Ports:     ports,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", s.path, err)
	}
	return nil
}

// Load reads the baseline snapshot from disk.
// Returns ErrNotFound if no baseline has been saved yet.
func (s *Store) Load() (*Snapshot, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("baseline: read %s: %w", s.path, err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("baseline: corrupt file %s: %w", s.path, err)
	}
	return &snap, nil
}

// ErrNotFound is returned when no baseline file exists on disk.
var ErrNotFound = fmt.Errorf("baseline: no baseline found")
