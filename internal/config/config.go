package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	ScanInterval time.Duration `yaml:"scan_interval"`
	Ports        []int         `yaml:"ports"`
	AlertFormat  string        `yaml:"alert_format"`
	LogFile      string        `yaml:"log_file"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		ScanInterval: 30 * time.Second,
		Ports:        []int{},
		AlertFormat:  "text",
		LogFile:      "",
	}
}

// Load reads a YAML config file from the given path and returns a Config.
// Fields not present in the file retain their default values.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: reading file %q: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config: parsing yaml: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	return cfg, nil
}

// validate checks that Config fields are within acceptable ranges.
func (c *Config) validate() error {
	if c.ScanInterval < time.Second {
		return fmt.Errorf("scan_interval must be at least 1s, got %s", c.ScanInterval)
	}
	if c.AlertFormat != "text" && c.AlertFormat != "json" {
		return fmt.Errorf("alert_format must be \"text\" or \"json\", got %q", c.AlertFormat)
	}
	for _, p := range c.Ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("port %d is out of valid range (1-65535)", p)
		}
	}
	return nil
}
