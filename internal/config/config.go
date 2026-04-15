// Package config loads and validates portwatch runtime configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all runtime settings for portwatch.
type Config struct {
	// ScanInterval is how often the port scanner runs.
	ScanInterval time.Duration `yaml:"scan_interval"`

	// Ports lists the TCP ports to monitor. An empty list means all ports
	// in the range 1–65535 are scanned.
	Ports []int `yaml:"ports"`

	// AlertFormat controls how change notifications are rendered.
	// Accepted values: "text", "json".
	AlertFormat string `yaml:"alert_format"`

	// StateFile is the path used to persist scan snapshots.
	StateFile string `yaml:"state_file"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		ScanInterval: 30 * time.Second,
		AlertFormat:  "text",
		StateFile:    "/tmp/portwatch_state.json",
	}
}

// Load reads a YAML configuration file from path and merges it with
// DefaultConfig, then validates the result.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing config: %w", err)
	}

	if err := validate(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func validate(cfg Config) error {
	if cfg.ScanInterval < time.Second {
		return fmt.Errorf("scan_interval must be at least 1s, got %s", cfg.ScanInterval)
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
