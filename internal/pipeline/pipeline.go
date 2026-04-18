// Package pipeline wires together the scan, filter, diff, and alert stages
// into a single reusable processing pipeline.
package pipeline

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
	"github.com/user/portwatch/internal/state"
)

// Stage represents a single processing step in the pipeline.
type Stage string

const (
	StageScan    Stage = "scan"
	StageFilter  Stage = "filter"
	StageDiff    Stage = "diff"
	StageAlert   Stage = "alert"
)

// Result holds the outcome of a single pipeline run.
type Result struct {
	ScannedAt  time.Time
	Snapshot   *snapshot.Snapshot
	Opened     []scanner.PortInfo
	Closed     []scanner.PortInfo
	Duration   time.Duration
}

// Pipeline coordinates scanning, filtering, diffing, and alerting.
type Pipeline struct {
	scanner  *scanner.TCPScanner
	filter   *filter.Filter
	store    *state.Store
	notifier *alert.Notifier
	metrics  *metrics.Metrics
	ports    []int
}

// Config holds the dependencies required to build a Pipeline.
type Config struct {
	Scanner  *scanner.TCPScanner
	Filter   *filter.Filter
	Store    *state.Store
	Notifier *alert.Notifier
	Metrics  *metrics.Metrics
	Ports    []int
}

// New constructs a Pipeline from the provided Config.
func New(cfg Config) (*Pipeline, error) {
	if cfg.Scanner == nil {
		return nil, fmt.Errorf("pipeline: scanner is required")
	}
	if cfg.Store == nil {
		return nil, fmt.Errorf("pipeline: state store is required")
	}
	return &Pipeline{
		scanner:  cfg.Scanner,
		filter:   cfg.Filter,
		store:    cfg.Store,
		notifier: cfg.Notifier,
		metrics:  cfg.Metrics,
		ports:    cfg.Ports,
	}, nil
}

// Run executes one full scan-filter-diff-alert cycle.
// It returns a Result summarising what changed, or an error if scanning failed.
func (p *Pipeline) Run(ctx context.Context) (*Result, error) {
	start := time.Now()

	// 1. Scan
	found, err := p.scanner.Scan(ctx, p.ports)
	if err != nil {
		return nil, fmt.Errorf("pipeline scan: %w", err)
	}
	if p.metrics != nil {
		p.metrics.RecordScan()
	}

	// 2. Filter
	if p.filter != nil {
		filtered := found[:0]
		for _, pi := range found {
			if !p.filter.Excluded(pi.Port) {
				filtered = append(filtered, pi)
			}
		}
		found = filtered
	}

	// 3. Build snapshot
	snap := snapshot.New(found)

	// 4. Diff against previous state
	prev, _ := p.store.Load()
	opened, closed := diff(prev, found)

	// 5. Persist new state
	if err := p.store.Save(found); err != nil {
		log.Printf("pipeline: failed to save state: %v", err)
	}

	// 6. Alert
	if p.notifier != nil {
		for _, pi := range opened {
			if alertErr := p.notifier.NotifyOpened(pi); alertErr != nil {
				log.Printf("pipeline: alert opened port %d: %v", pi.Port, alertErr)
			}
			if p.metrics != nil {
				p.metrics.RecordAlert()
			}
		}
		for _, pi := range closed {
			if alertErr := p.notifier.NotifyClosed(pi); alertErr != nil {
				log.Printf("pipeline: alert closed port %d: %v", pi.Port, alertErr)
			}
			if p.metrics != nil {
				p.metrics.RecordAlert()
			}
		}
	}

	return &Result{
		ScannedAt: start,
		Snapshot:  snap,
		Opened:    opened,
		Closed:    closed,
		Duration:  time.Since(start),
	}, nil
}

// diff computes which ports were opened or closed between two scans.
func diff(prev, curr []scanner.PortInfo) (opened, closed []scanner.PortInfo) {
	prevMap := make(map[int]scanner.PortInfo, len(prev))
	for _, p := range prev {
		prevMap[p.Port] = p
	}
	currMap := make(map[int]scanner.PortInfo, len(curr))
	for _, p := range curr {
		currMap[p.Port] = p
	}
	for port, pi := range currMap {
		if _, ok := prevMap[port]; !ok {
			opened = append(opened, pi)
		}
	}
	for port, pi := range prevMap {
		if _, ok := currMap[port]; !ok {
			closed = append(closed, pi)
		}
	}
	return
}
