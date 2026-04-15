package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestMain_VersionFlag builds the binary and checks that --version exits 0
// and prints the expected prefix. Skipped when the go tool is unavailable.
func TestMain_VersionFlag(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go binary not found in PATH")
	}

	tmp := t.TempDir()
	bin := tmp + "/portwatch"

	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	cmd := exec.Command(bin, "--version")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("--version exited with error: %v", err)
	}

	if !strings.HasPrefix(string(out), "portwatch ") {
		t.Errorf("unexpected version output: %q", string(out))
	}
}

// TestMain_InvalidConfig ensures the process exits non-zero when given a
// config file that contains invalid TOML.
func TestMain_InvalidConfig(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go binary not found in PATH")
	}

	tmp := t.TempDir()
	bin := tmp + "/portwatch"

	build := exec.Command("go", "build", "-o", bin, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	badCfg := tmp + "/bad.toml"
	if err := os.WriteFile(badCfg, []byte("!!not toml!!"), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin, "--config", badCfg)
	if err := cmd.Run(); err == nil {
		t.Error("expected non-zero exit for invalid config, got nil")
	}
}
