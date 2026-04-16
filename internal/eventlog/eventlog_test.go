package eventlog_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/eventlog"
	"github.com/user/portwatch/internal/scanner"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "events.jsonl")
}

func samplePort() scanner.PortInfo {
	return scanner.PortInfo{Port: 8080, Proto: "tcp", State: "open"}
}

func TestAppend_CreatesFile(t *testing.T) {
	s := eventlog.NewStore(tempPath(t))
	if err := s.Append("opened", samplePort()); err != nil {
		t.Fatalf("Append: %v", err)
	}
}

func TestReadAll_EmptyWhenMissing(t *testing.T) {
	s := eventlog.NewStore(tempPath(t))
	events, err := s.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

func TestAppendAndReadAll_RoundTrip(t *testing.T) {
	s := eventlog.NewStore(tempPath(t))
	p := samplePort()

	_ = s.Append("opened", p)
	_ = s.Append("closed", p)

	events, err := s.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Kind != "opened" || events[1].Kind != "closed" {
		t.Errorf("unexpected kinds: %v", events)
	}
	if events[0].Port.Port != 8080 {
		t.Errorf("unexpected port: %v", events[0].Port)
	}
}

func TestAppend_TimestampIsRecent(t *testing.T) {
	s := eventlog.NewStore(tempPath(t))
	before := time.Now().UTC()
	_ = s.Append("opened", samplePort())
	after := time.Now().UTC()

	events, _ := s.ReadAll()
	ts := events[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v out of range [%v, %v]", ts, before, after)
	}
}
