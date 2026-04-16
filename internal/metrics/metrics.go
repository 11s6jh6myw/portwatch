// Package metrics tracks runtime statistics for the portwatch daemon.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of daemon metrics.
type Snapshot struct {
	ScansTotal    int64         `json:"scans_total"`
	AlertsTotal   int64         `json:"alerts_total"`
	OpenPorts     int           `json:"open_ports"`
	LastScanAt    time.Time     `json:"last_scan_at"`
	LastScanDur   time.Duration `json:"last_scan_duration_ms"`
	UptimeSeconds float64       `json:"uptime_seconds"`
}

// Collector accumulates daemon metrics.
type Collector struct {
	mu          sync.Mutex
	scansTotal  int64
	alertsTotal int64
	openPorts   int
	lastScanAt  time.Time
	lastScanDur time.Duration
	startedAt   time.Time
}

// New returns a new Collector initialised with the current time.
func New() *Collector {
	return &Collector{startedAt: time.Now()}
}

// RecordScan records a completed scan.
func (c *Collector) RecordScan(dur time.Duration, openPorts int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scansTotal++
	c.lastScanAt = time.Now()
	c.lastScanDur = dur
	c.openPorts = openPorts
}

// RecordAlert increments the alert counter.
func (c *Collector) RecordAlert() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alertsTotal++
}

// Snapshot returns a copy of current metrics.
func (c *Collector) Snapshot() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Snapshot{
		ScansTotal:    c.scansTotal,
		AlertsTotal:   c.alertsTotal,
		OpenPorts:     c.openPorts,
		LastScanAt:    c.lastScanAt,
		LastScanDur:   c.lastScanDur,
		UptimeSeconds: time.Since(c.startedAt).Seconds(),
	}
}
