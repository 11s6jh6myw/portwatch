// Package schedule provides adaptive scan interval adjustment based on
// recent port activity. When changes are detected frequently the interval
// shrinks; when the environment is quiet it grows back toward the maximum.
package schedule

import (
	"sync"
	"time"
)

// Config holds tuning parameters for the adaptive scheduler.
type Config struct {
	MinInterval time.Duration
	MaxInterval time.Duration
	// StepDown is the factor applied to the interval on activity (< 1).
	StepDown float64
	// StepUp is the factor applied to the interval on quiet scans (> 1).
	StepUp float64
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		MinInterval: 5 * time.Second,
		MaxInterval: 60 * time.Second,
		StepDown:    0.5,
		StepUp:      1.25,
	}
}

// Scheduler tracks the current scan interval and adjusts it adaptively.
type Scheduler struct {
	cfg     Config
	mu      sync.Mutex
	current time.Duration
}

// New creates a Scheduler starting at the maximum interval.
func New(cfg Config) *Scheduler {
	return &Scheduler{cfg: cfg, current: cfg.MaxInterval}
}

// Current returns the current scan interval.
func (s *Scheduler) Current() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current
}

// RecordActivity signals that changes were detected; the interval shrinks.
func (s *Scheduler) RecordActivity() {
	s.mu.Lock()
	defer s.mu.Unlock()
	next := time.Duration(float64(s.current) * s.cfg.StepDown)
	if next < s.cfg.MinInterval {
		next = s.cfg.MinInterval
	}
	s.current = next
}

// RecordQuiet signals that no changes were detected; the interval grows.
func (s *Scheduler) RecordQuiet() {
	s.mu.Lock()
	defer s.mu.Unlock()
	next := time.Duration(float64(s.current) * s.cfg.StepUp)
	if next > s.cfg.MaxInterval {
		next = s.cfg.MaxInterval
	}
	s.current = next
}
