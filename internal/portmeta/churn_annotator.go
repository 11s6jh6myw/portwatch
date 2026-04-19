package portmeta

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

const (
	churnLevelKey  = "churn.level"
	churnCountKey  = "churn.count"
	churnWindowKey = "churn.window"
)

// ChurnAnnotator attaches churn metadata to ports based on recent events.
type ChurnAnnotator struct {
	events map[int][]ChurnEvent
	window time.Duration
}

// NewChurnAnnotator creates an annotator using the provided event history
// and observation window.
func NewChurnAnnotator(events map[int][]ChurnEvent, window time.Duration) *ChurnAnnotator {
	return &ChurnAnnotator{events: events, window: window}
}

// Annotate adds churn level, count, and window metadata to each port.
func (a *ChurnAnnotator) Annotate(ports []scanner.PortInfo) []scanner.PortInfo {
	out := make([]scanner.PortInfo, len(ports))
	for i, p := range ports {
		evs := a.events[p.Port]
		level := ClassifyChurn(evs, a.window)
		count := countWithinWindow(evs, a.window)
		if p.Meta == nil {
			p.Meta = make(map[string]string)
		}
		p.Meta[churnLevelKey] = level.String()
		p.Meta[churnCountKey] = fmt.Sprintf("%d", count)
		p.Meta[churnWindowKey] = a.window.String()
		out[i] = p
	}
	return out
}

// FilterByMaxChurn returns only ports whose churn level is at or below max.
func FilterByMaxChurn(ports []scanner.PortInfo, max ChurnLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if p.Meta == nil {
			out = append(out, p)
			continue
		}
		level := parseChurnLevel(p.Meta[churnLevelKey])
		if level <= max {
			out = append(out, p)
		}
	}
	return out
}

func countWithinWindow(events []ChurnEvent, window time.Duration) int {
	cutoff := time.Now().Add(-window)
	n := 0
	for _, e := range events {
		if e.At.After(cutoff) {
			n++
		}
	}
	return n
}

func parseChurnLevel(s string) ChurnLevel {
	switch s {
	case "low":
		return ChurnLow
	case "moderate":
		return ChurnModerate
	case "high":
		return ChurnHigh
	default:
		return ChurnNone
	}
}
