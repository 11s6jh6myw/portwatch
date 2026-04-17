package fingerprint_test

import (
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/scanner"
)

func TestCache_FirstEntryAlwaysChanged(t *testing.T) {
	c := fingerprint.NewCache()
	ports := []scanner.PortInfo{makePort(80, "tcp", "open")}
	fp := fingerprint.Compute(ports)
	if !c.Changed("host1", fp) {
		t.Fatal("first entry should always be reported as changed")
	}
}

func TestCache_SameFingerprintNotChanged(t *testing.T) {
	c := fingerprint.NewCache()
	ports := []scanner.PortInfo{makePort(80, "tcp", "open")}
	fp := fingerprint.Compute(ports)
	c.Changed("host1", fp)
	if c.Changed("host1", fp) {
		t.Fatal("same fingerprint should not be reported as changed")
	}
}

func TestCache_DifferentFingerprintChanged(t *testing.T) {
	c := fingerprint.NewCache()
	fp1 := fingerprint.Compute([]scanner.PortInfo{makePort(80, "tcp", "open")})
	fp2 := fingerprint.Compute([]scanner.PortInfo{makePort(443, "tcp", "open")})
	c.Changed("host1", fp1)
	if !c.Changed("host1", fp2) {
		t.Fatal("different fingerprint should be reported as changed")
	}
}

func TestCache_Get_ReturnsStored(t *testing.T) {
	c := fingerprint.NewCache()
	fp := fingerprint.Compute([]scanner.PortInfo{makePort(22, "tcp", "open")})
	c.Changed("h", fp)
	got, ok := c.Get("h")
	if !ok || got != fp {
		t.Fatalf("expected %v got %v (ok=%v)", fp, got, ok)
	}
}

func TestCache_Invalidate_RemovesEntry(t *testing.T) {
	c := fingerprint.NewCache()
	fp := fingerprint.Compute([]scanner.PortInfo{makePort(22, "tcp", "open")})
	c.Changed("h", fp)
	c.Invalidate("h")
	_, ok := c.Get("h")
	if ok {
		t.Fatal("entry should have been removed after Invalidate")
	}
}

func TestCache_IndependentKeys(t *testing.T) {
	c := fingerprint.NewCache()
	fp1 := fingerprint.Compute([]scanner.PortInfo{makePort(80, "tcp", "open")})
	fp2 := fingerprint.Compute([]scanner.PortInfo{makePort(443, "tcp", "open")})
	c.Changed("a", fp1)
	c.Changed("b", fp2)
	if c.Changed("a", fp1) {
		t.Fatal("key a should not be changed")
	}
	if c.Changed("b", fp2) {
		t.Fatal("key b should not be changed")
	}
}
