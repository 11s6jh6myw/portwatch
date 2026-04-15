package history_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/scanner"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func samplePort() scanner.PortInfo {
	return scanner.PortInfo{Port: 8080, Protocol: "tcp", State: "open"}
}

func TestStore_RecordAndRetrieve(t *testing.T) {
	s, err := history.NewStore(tempPath(t))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	if err := s.Record(history.EventOpened, samplePort()); err != nil {
		t.Fatalf("Record: %v", err)
	}

	events := s.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Type != history.EventOpened {
		t.Errorf("expected type %q, got %q", history.EventOpened, events[0].Type)
	}
	if events[0].Port.Port != 8080 {
		t.Errorf("expected port 8080, got %d", events[0].Port.Port)
	}
}

func TestStore_PersistsAcrossReloads(t *testing.T) {
	path := tempPath(t)

	s1, _ := history.NewStore(path)
	_ = s1.Record(history.EventOpened, samplePort())
	_ = s1.Record(history.EventClosed, samplePort())

	s2, err := history.NewStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(s2.Events()) != 2 {
		t.Errorf("expected 2 events after reload, got %d", len(s2.Events()))
	}
}

func TestStore_LoadMissingFile(t *testing.T) {
	s, err := history.NewStore(tempPath(t))
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(s.Events()) != 0 {
		t.Errorf("expected empty history, got %d events", len(s.Events()))
	}
}

func TestStore_LoadCorruptFile(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not valid json{"), 0o644)

	_, err := history.NewStore(path)
	if err == nil {
		t.Fatal("expected error loading corrupt file, got nil")
	}
}

func TestStore_EventsReturnsCopy(t *testing.T) {
	s, _ := history.NewStore(tempPath(t))
	_ = s.Record(history.EventOpened, samplePort())

	events := s.Events()
	events[0].Type = history.EventClosed

	original := s.Events()
	if original[0].Type != history.EventOpened {
		t.Error("Events() should return a copy, not a reference")
	}
}
