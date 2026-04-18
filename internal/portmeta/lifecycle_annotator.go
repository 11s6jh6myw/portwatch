package portmeta

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// LifecycleAnnotator attaches a lifecycle stage to each PortInfo.
type LifecycleAnnotator struct {
	entries      map[string]LifecycleEntry
	stableAfter  time.Duration
}

// NewLifecycleAnnotator creates an annotator with the given stability threshold.
func NewLifecycleAnnotator(stableAfter time.Duration) *LifecycleAnnotator {
	return &LifecycleAnnotator{
		entries:     make(map[string]LifecycleEntry),
		stableAfter: stableAfter,
	}
}

func lifecycleKey(p scanner.PortInfo) string {
	return p.Protocol + ":" + itoa(p.Port)
}

func itoa(n int) string {
	const digits = "0123456789"
	if n == 0 {
		return "0"
	}
	buf := [10]byte{}
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = digits[n%10]
		n /= 10
	}
	return string(buf[i:])
}

// RecordOpen updates the entry for an opened port.
func (a *LifecycleAnnotator) RecordOpen(p scanner.PortInfo, now time.Time) {
	k := lifecycleKey(p)
	e := a.entries[k]
	if e.FirstSeen.IsZero() {
		e.FirstSeen = now
	}
	e.Port = p.Port
	e.Protocol = p.Protocol
	e.OpenCount++
	e.LastSeen = now
	a.entries[k] = e
}

// RecordClose updates the entry for a closed port.
func (a *LifecycleAnnotator) RecordClose(p scanner.PortInfo, now time.Time) {
	k := lifecycleKey(p)
	e := a.entries[k]
	e.CloseCount++
	e.OpenCount = max0(e.OpenCount - 1)
	e.LastSeen = now
	a.entries[k] = e
}

// Stage returns the current lifecycle stage for the given port.
func (a *LifecycleAnnotator) Stage(p scanner.PortInfo, now time.Time) LifecycleStage {
	e, ok := a.entries[lifecycleKey(p)]
	if !ok {
		return StageNew
	}
	return Classify(e, now, a.stableAfter)
}

func max0(n int) int {
	if n < 0 {
		return 0
	}
	return n
}
