package portmeta

import "time"

// ChurnLevel describes how frequently a port's state has changed.
type ChurnLevel int

const (
	ChurnNone     ChurnLevel = iota // 0 changes
	ChurnLow                        // 1–2 changes
	ChurnModerate                   // 3–5 changes
	ChurnHigh                       // 6+ changes
)

func (c ChurnLevel) String() string {
	switch c {
	case ChurnNone:
		return "none"
	case ChurnLow:
		return "low"
	case ChurnModerate:
		return "moderate"
	case ChurnHigh:
		return "high"
	default:
		return "unknown"
	}
}

// ChurnEvent records a single open/close state transition.
type ChurnEvent struct {
	At   time.Time
	Kind string // "opened" or "closed"
}

// ClassifyChurn returns a ChurnLevel based on the number of transitions
// observed within the given window.
func ClassifyChurn(events []ChurnEvent, window time.Duration) ChurnLevel {
	if len(events) == 0 {
		return ChurnNone
	}
	cutoff := time.Now().Add(-window)
	count := 0
	for _, e := range events {
		if e.At.After(cutoff) {
			count++
		}
	}
	switch {
	case count == 0:
		return ChurnNone
	case count <= 2:
		return ChurnLow
	case count <= 5:
		return ChurnModerate
	default:
		return ChurnHigh
	}
}
