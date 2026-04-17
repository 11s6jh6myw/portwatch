package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/audit"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.jsonl")
}

func sampleEntry() audit.Entry {
	return audit.Entry{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Hostname:  "testhost",
		Kind:      "opened",
		Port:      8080,
		Protocol:  "tcp",
		Action:    "alerted",
	}
}

func TestRecord_CreatesFile(t *testing.T) {
	path := tempPath(t)
	s := audit.NewStore(path)
	if err := s.Record(sampleEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestReadAll_EmptyWhenMissing(t *testing.T) {
	s := audit.NewStore(tempPath(t))
	entries, err := s.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestRecord_AndReadAll_RoundTrip(t *testing.T) {
	s := audit.NewStore(tempPath(t))
	e1 := sampleEntry()
	e2 := sampleEntry()
	e2.Kind = "closed"
	e2.Action = "suppressed"
	e2.Port = 9090

	if err := s.Record(e1); err != nil {
		t.Fatal(err)
	}
	if err := s.Record(e2); err != nil {
		t.Fatal(err)
	}

	entries, err := s.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Port != 8080 || entries[0].Kind != "opened" {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
	if entries[1].Port != 9090 || entries[1].Action != "suppressed" {
		t.Errorf("unexpected second entry: %+v", entries[1])
	}
}

func TestRecord_SetsTimestampIfZero(t *testing.T) {
	s := audit.NewStore(tempPath(t))
	e := sampleEntry()
	e.Timestamp = time.Time{}
	if err := s.Record(e); err != nil {
		t.Fatal(err)
	}
	entries, _ := s.ReadAll()
	if entries[0].Timestamp.IsZero() {
		t.Error("expected timestamp to be set")
	}
}
