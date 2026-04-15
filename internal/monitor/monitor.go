package monitor

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// AlertFunc is called when port changes are detected.
type AlertFunc func(diff scanner.DiffResult)

// Monitor periodically scans ports and alerts on changes.
type Monitor struct {
	scanner  *scanner.TCPScanner
	ports    []int
	interval time.Duration
	alert    AlertFunc
	stopCh   chan struct{}
}

// Config holds configuration for the Monitor.
type Config struct {
	Ports    []int
	Interval time.Duration
	Alert    AlertFunc
}

// New creates a new Monitor with the given config.
func New(cfg Config) *Monitor {
	return &Monitor{
		scanner:  scanner.NewTCPScanner(500 * time.Millisecond),
		ports:    cfg.Ports,
		interval: cfg.Interval,
		alert:    cfg.Alert,
		stopCh:   make(chan struct{}),
	}
}

// Start begins monitoring in a blocking loop until Stop is called.
func (m *Monitor) Start() {
	log.Printf("portwatch: starting monitor, interval=%s ports=%v", m.interval, m.ports)

	previous := m.scan()
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			current := m.scan()
			diff := scanner.Diff(previous, current)
			if diff.HasChanges() {
				m.alert(diff)
			}
			previous = current
		case <-m.stopCh:
			log.Println("portwatch: monitor stopped")
			return
		}
	}
}

// Stop signals the monitor to stop.
func (m *Monitor) Stop() {
	close(m.stopCh)
}

func (m *Monitor) scan() []scanner.PortInfo {
	results, err := m.scanner.Scan(m.ports)
	if err != nil {
		log.Printf("portwatch: scan error: %v", err)
		return nil
	}
	return results
}
