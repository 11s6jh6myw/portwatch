package snapshot_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makePort(port int) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Proto: "tcp", State: "open"}
}

func TestNew_SetsTimestampAndHostname(t *testing.T) {
	before := time.Now().UTC()
	snap := snapshot.New("localhost", []scanner.PortInfo{makePort(80)})
	after := time.Now().UTC()

	if snap.Hostname != "localhost" {
		t.Errorf("expected hostname 'localhost', got %q", snap.Hostname)
	}
	if snap.Timestamp.Before(before) || snap.Timestamp.After(after) {
		t.Error("timestamp out of expected range")
	}
	if len(snap.Ports) != 1 {
		t.Errorf("expected 1 port, got %d", len(snap.Ports))
	}
}

func TestSummary_ContainsExpectedFields(t *testing.T) {
	snap := snapshot.New("myhost", []scanner.PortInfo{makePort(22), makePort(443)})
	summary := snap.Summary()

	if !strings.Contains(summary, "myhost") {
		t.Errorf("summary missing hostname: %q", summary)
	}
	if !strings.Contains(summary, "2 open port(s)") {
		t.Errorf("summary missing port count: %q", summary)
	}
}

func TestPortSet_ReturnsMapByPort(t *testing.T) {
	ports := []scanner.PortInfo{makePort(22), makePort(80), makePort(443)}
	snap := snapshot.New("localhost", ports)
	set := snap.PortSet()

	for _, p := range ports {
		if _, ok := set[p.Port]; !ok {
			t.Errorf("port %d missing from PortSet", p.Port)
		}
	}
}

func TestEqual_SamePortsReturnsTrue(t *testing.T) {
	a := snapshot.New("h", []scanner.PortInfo{makePort(22), makePort(80)})
	b := snapshot.New("h", []scanner.PortInfo{makePort(80), makePort(22)})
	if !a.Equal(b) {
		t.Error("expected snapshots to be equal")
	}
}

func TestEqual_DifferentPortsReturnsFalse(t *testing.T) {
	a := snapshot.New("h", []scanner.PortInfo{makePort(22)})
	b := snapshot.New("h", []scanner.PortInfo{makePort(22), makePort(80)})
	if a.Equal(b) {
		t.Error("expected snapshots to be unequal")
	}
}

func TestEqual_EmptySnapshots(t *testing.T) {
	a := snapshot.New("h", nil)
	b := snapshot.New("h", []scanner.PortInfo{})
	if !a.Equal(b) {
		t.Error("expected two empty snapshots to be equal")
	}
}
