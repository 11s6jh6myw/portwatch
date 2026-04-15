// Package state provides persistent storage for port scan snapshots,
// allowing portwatch to detect changes across restarts.
package state

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot holds a recorded set of open ports at a point in time.
type Snapshot struct {
	Timestamp time.Time          `json:"timestamp"`
	Ports     []scanner.PortInfo `json:"ports"`
}

// Store manages reading and writing snapshots to disk.
type Store struct {
	path string
}

// NewStore creates a Store that persists data at the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the given ports as the latest snapshot to disk.
func (s *Store) Save(ports []scanner.PortInfo) error {
	snap := Snapshot{
		Timestamp: time.Now(),
		Ports:     ports,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Load reads the latest snapshot from disk.
// Returns an empty Snapshot and no error if the file does not exist yet.
func (s *Store) Load() (Snapshot, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return Snapshot{}, nil
		}
		return Snapshot{}, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}
