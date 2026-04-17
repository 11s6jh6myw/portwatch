package rollup_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/joshbeard/portwatch/internal/rollup"
	"github.com/joshbeard/portwatch/internal/scanner"
)

func makeSummary() rollup.Summary {
	return rollup.Summary{
		Opened: []scanner.PortInfo{{Port: 80, Proto: "tcp"}, {Port: 443, Proto: "tcp"}},
		Closed: []scanner.PortInfo{{Port: 8080, Proto: "tcp"}},
		At:     time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}
}

func TestFormat_Text(t *testing.T) {
	s := makeSummary()
	out, err := rollup.Format(s, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "opened (2)") {
		t.Errorf("missing opened count: %q", out)
	}
	if !strings.Contains(out, "closed (1)") {
		t.Errorf("missing closed count: %q", out)
	}
	if !strings.Contains(out, "80/tcp") {
		t.Errorf("missing port 80: %q", out)
	}
}

func TestFormat_JSON(t *testing.T) {
	s := makeSummary()
	out, err := rollup.Format(s, "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var v struct {
		Opened []int  `json:"opened"`
		Closed []int  `json:"closed"`
		At     string `json:"at"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(v.Opened) != 2 {
		t.Errorf("expected 2 opened, got %d", len(v.Opened))
	}
	if len(v.Closed) != 1 {
		t.Errorf("expected 1 closed, got %d", len(v.Closed))
	}
}

func TestFormat_DefaultsToText(t *testing.T) {
	s := makeSummary()
	out, err := rollup.Format(s, "unknown")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[rollup]") {
		t.Errorf("expected text format, got: %q", out)
	}
}
