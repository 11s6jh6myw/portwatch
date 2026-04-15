package config

import (
	"os"
	"testing"
	"time"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.ScanInterval != 30*time.Second {
		t.Errorf("expected 30s scan interval, got %s", cfg.ScanInterval)
	}
	if cfg.AlertFormat != "text" {
		t.Errorf("expected alert_format=text, got %q", cfg.AlertFormat)
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, `
scan_interval: 10s
ports: [80, 443, 8080]
alert_format: json
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval != 10*time.Second {
		t.Errorf("expected 10s, got %s", cfg.ScanInterval)
	}
	if len(cfg.Ports) != 3 {
		t.Errorf("expected 3 ports, got %d", len(cfg.Ports))
	}
	if cfg.AlertFormat != "json" {
		t.Errorf("expected json, got %q", cfg.AlertFormat)
	}
}

func TestLoad_InvalidScanInterval(t *testing.T) {
	path := writeTemp(t, "scan_interval: 500ms\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for scan_interval < 1s")
	}
}

func TestLoad_InvalidAlertFormat(t *testing.T) {
	path := writeTemp(t, "alert_format: xml\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for unknown alert_format")
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	path := writeTemp(t, "ports: [0, 80]\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for port 0")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/portwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
