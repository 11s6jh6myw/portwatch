package metrics_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestNew_InitialisesZeroCounters(t *testing.T) {
	c := metrics.New()
	s := c.Snapshot()
	if s.ScansTotal != 0 || s.AlertsTotal != 0 {
		t.Fatalf("expected zero counters, got scans=%d alerts=%d", s.ScansTotal, s.AlertsTotal)
	}
}

func TestRecordScan_IncrementsCounter(t *testing.T) {
	c := metrics.New()
	c.RecordScan(10*time.Millisecond, 5)
	c.RecordScan(20*time.Millisecond, 3)
	s := c.Snapshot()
	if s.ScansTotal != 2 {
		t.Fatalf("expected 2 scans, got %d", s.ScansTotal)
	}
	if s.OpenPorts != 3 {
		t.Fatalf("expected 3 open ports, got %d", s.OpenPorts)
	}
	if s.LastScanDur != 20*time.Millisecond {
		t.Fatalf("unexpected last scan duration: %v", s.LastScanDur)
	}
}

func TestRecordAlert_IncrementsCounter(t *testing.T) {
	c := metrics.New()
	c.RecordAlert()
	c.RecordAlert()
	c.RecordAlert()
	s := c.Snapshot()
	if s.AlertsTotal != 3 {
		t.Fatalf("expected 3 alerts, got %d", s.AlertsTotal)
	}
}

func TestSnapshot_UptimeIncreases(t *testing.T) {
	c := metrics.New()
	time.Sleep(10 * time.Millisecond)
	s := c.Snapshot()
	if s.UptimeSeconds <= 0 {
		t.Fatalf("expected positive uptime, got %f", s.UptimeSeconds)
	}
}

func TestSnapshot_LastScanAt_Updated(t *testing.T) {
	c := metrics.New()
	before := time.Now()
	c.RecordScan(5*time.Millisecond, 1)
	s := c.Snapshot()
	if s.LastScanAt.Before(before) {
		t.Fatalf("LastScanAt not updated")
	}
}
