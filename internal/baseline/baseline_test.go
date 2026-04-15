package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/scanner"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func samplePorts() []scanner.PortInfo {
	return []scanner.PortInfo{
		{Port: 22, Proto: "tcp", State: "open"},
		{Port: 443, Proto: "tcp", State: "open"},
	}
}

func TestStore_SaveAndLoad(t *testing.T) {
	store := baseline.NewStore(tempPath(t))
	ports := samplePorts()

	if err := store.Save(ports); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(snap.Ports) != len(ports) {
		t.Errorf("got %d ports, want %d", len(snap.Ports), len(ports))
	}
	if snap.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestStore_LoadMissingFile(t *testing.T) {
	store := baseline.NewStore(tempPath(t))
	_, err := store.Load()
	if err != baseline.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestStore_LoadCorruptFile(t *testing.T) {
	path := tempPath(t)
	if err := os.WriteFile(path, []byte("not-json{"), 0o644); err != nil {
		t.Fatal(err)
	}
	store := baseline.NewStore(path)
	_, err := store.Load()
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}

func TestStore_OverwritesPreviousSave(t *testing.T) {
	store := baseline.NewStore(tempPath(t))

	if err := store.Save(samplePorts()); err != nil {
		t.Fatal(err)
	}
	newPorts := []scanner.PortInfo{{Port: 8080, Proto: "tcp", State: "open"}}
	if err := store.Save(newPorts); err != nil {
		t.Fatal(err)
	}

	snap, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(snap.Ports) != 1 || snap.Ports[0].Port != 8080 {
		t.Errorf("unexpected ports after overwrite: %+v", snap.Ports)
	}
}
