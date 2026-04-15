package state_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func TestStore_SaveAndLoad(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 80, Protocol: "tcp"},
		{Port: 443, Protocol: "tcp"},
	}
	s := state.NewStore(tempPath(t))

	if err := s.Save(ports); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(snap.Ports) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(snap.Ports))
	}
	for i, p := range ports {
		if snap.Ports[i].Port != p.Port {
			t.Errorf("port[%d]: expected %d, got %d", i, p.Port, snap.Ports[i].Port)
		}
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if snap.Timestamp.After(time.Now().Add(time.Second)) {
		t.Error("timestamp is in the future")
	}
}

func TestStore_LoadMissingFile(t *testing.T) {
	s := state.NewStore(tempPath(t))
	snap, err := s.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if snap.Ports != nil {
		t.Errorf("expected nil ports for empty snapshot, got %v", snap.Ports)
	}
}

func TestStore_LoadCorruptFile(t *testing.T) {
	p := tempPath(t)
	if err := os.WriteFile(p, []byte("not json{"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := state.NewStore(p)
	_, err := s.Load()
	if err == nil {
		t.Fatal("expected error for corrupt file, got nil")
	}
}

func TestStore_OverwritesPreviousSnapshot(t *testing.T) {
	s := state.NewStore(tempPath(t))

	first := []scanner.PortInfo{{Port: 22, Protocol: "tcp"}}
	if err := s.Save(first); err != nil {
		t.Fatal(err)
	}

	second := []scanner.PortInfo{{Port: 8080, Protocol: "tcp"}, {Port: 9090, Protocol: "tcp"}}
	if err := s.Save(second); err != nil {
		t.Fatal(err)
	}

	snap, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(snap.Ports) != 2 {
		t.Fatalf("expected 2 ports after overwrite, got %d", len(snap.Ports))
	}
}
