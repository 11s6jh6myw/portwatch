// Package snapshot provides types and utilities for capturing and comparing
// point-in-time views of open ports on the system.
package snapshot

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot represents a point-in-time capture of open ports.
type Snapshot struct {
	Timestamp time.Time          `json:"timestamp"`
	Ports     []scanner.PortInfo `json:"ports"`
	Hostname  string             `json:"hostname"`
}

// New creates a new Snapshot with the current timestamp.
func New(hostname string, ports []scanner.PortInfo) *Snapshot {
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
		Hostname:  hostname,
	}
}

// Summary returns a human-readable summary of the snapshot.
func (s *Snapshot) Summary() string {
	return fmt.Sprintf("snapshot at %s: %d open port(s) on %s",
		s.Timestamp.Format(time.RFC3339),
		len(s.Ports),
		s.Hostname,
	)
}

// PortSet returns a map of port numbers to PortInfo for fast lookup.
func (s *Snapshot) PortSet() map[int]scanner.PortInfo {
	m := make(map[int]scanner.PortInfo, len(s.Ports))
	for _, p := range s.Ports {
		m[p.Port] = p
	}
	return m
}

// Equal reports whether two snapshots contain the same set of open ports.
func (s *Snapshot) Equal(other *Snapshot) bool {
	if len(s.Ports) != len(other.Ports) {
		return false
	}
	set := s.PortSet()
	for _, p := range other.Ports {
		if _, ok := set[p.Port]; !ok {
			return false
		}
	}
	return true
}
