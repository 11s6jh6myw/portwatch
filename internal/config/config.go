// Package config handles loading and validating portwatch configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the full portwatch runtime configuration.
type Config struct {
	// ScanInterval is how often the port scanner runs.
	ScanInterval time.Duration `yaml:"scan_interval"`

	// Ports is the list of TCP ports to monitor.
	Ports []uint16 `yaml:"ports"`

	// AlertFormat controls alert output: "text" or "json".
	AlertFormat string `yaml:"alert_format"`

	// AlertOutput is the file path for alerts; empty means stdout.
	AlertOutput string `yaml:"alert_output"`

	// StateFile is the path where port snapshots are persisted.
	StateFile string `yaml:"state_file"`

	// ExcludePorts lists ports or ranges to ignore (e.g. "22", "8000-9000").
	ExcludePorts []string `yaml:"exclude_ports"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		ScanInterval: 30 * time.Second,
		AlertFormat:  "text",
		StateFile:    "/tmp/portwatch.state.json",
		Ports:        []uint16{},
		ExcludePorts: []string{},
	}
}

// Load reads a YAML config file from path, merges it over defaults,
// and validates the result.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config: parse %q: %w", path, err)
	}

	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("config: invalid: %w", err)
	}

	return cfg, nil
}

func validate(cfg *Config) error {
	if cfg.ScanInterval < time.Second {
		return errors.New("scan_interval must be at least 1s")
	}
	switch cfg.AlertFormat {
	case "text", "json":
	default:
		return fmt.Errorf("alert_format must be \"text\" or \"json\", got %q", cfg.AlertFormat)
	}
	if cfg.StateFile == "" {
		return errors.New("state_file must not be empty")
	}
	return nil
}
