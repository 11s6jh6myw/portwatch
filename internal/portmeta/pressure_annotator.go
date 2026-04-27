package portmeta

import (
	"strconv"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

const (
	metaKeyPressure       = "pressure"
	metaKeyPressureCount  = "pressure_count"
	metaKeyPressureWindow = "pressure_window_s"
)

// NewPressureAnnotator returns an annotator that records the pressure level
// for each port into its Meta map. It counts events from the provided
// eventTimes map (keyed by port number) within window.
func NewPressureAnnotator(eventTimes map[int][]time.Time, window time.Duration) func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			events := eventTimes[p.Port]
			level := PressureFor(p, events, window)

			if p.Meta == nil {
				p.Meta = make(map[string]string)
			}
			p.Meta[metaKeyPressure] = level.String()
			p.Meta[metaKeyPressureCount] = strconv.Itoa(countWithinPressureWindow(events, window))
			p.Meta[metaKeyPressureWindow] = strconv.FormatInt(int64(window.Seconds()), 10)
			out[i] = p
		}
		return out
	}
}

func countWithinPressureWindow(events []time.Time, window time.Duration) int {
	if window <= 0 {
		window = time.Hour
	}
	cutoff := time.Now().Add(-window)
	n := 0
	for _, t := range events {
		if t.After(cutoff) {
			n++
		}
	}
	return n
}

// FilterByMaxPressure returns only ports whose recorded pressure level is at
// or below max. Ports without pressure metadata are passed through unchanged.
func FilterByMaxPressure(ports []scanner.PortInfo, max PressureLevel) []scanner.PortInfo {
	out := ports[:0:0]
	for _, p := range ports {
		raw, ok := p.Meta[metaKeyPressure]
		if !ok {
			out = append(out, p)
			continue
		}
		if parsePressure(raw) <= max {
			out = append(out, p)
		}
	}
	return out
}

func parsePressure(s string) PressureLevel {
	switch s {
	case "low":
		return PressureLow
	case "moderate":
		return PressureModerate
	case "high":
		return PressureHigh
	case "critical":
		return PressureCritical
	default:
		return PressureNone
	}
}
