package portmeta

import (
	"strconv"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

const (
	metaVelocity       = "velocity"
	metaVelocityCount  = "velocity_count"
	defaultVelocityWin = 30 * time.Minute
)

// NewVelocityAnnotator returns an annotator that computes and attaches
// velocity metadata to each PortInfo based on recent event timestamps
// stored in the port's existing metadata.
func NewVelocityAnnotator(window time.Duration) func([]scanner.PortInfo) []scanner.PortInfo {
	if window <= 0 {
		window = defaultVelocityWin
	}
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			events := extractEventTimes(p)
			level := VelocityFor(events, window)
			if p.Meta == nil {
				p.Meta = map[string]string{}
			}
			p.Meta[metaVelocity] = level.String()
			p.Meta[metaVelocityCount] = strconv.Itoa(len(events))
			out[i] = p
		}
		return out
	}
}

// FilterByMaxVelocity returns only ports whose velocity is at or below max.
func FilterByMaxVelocity(ports []scanner.PortInfo, max VelocityLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if p.Meta == nil {
			out = append(out, p)
			continue
		}
		level := parseVelocity(p.Meta[metaVelocity])
		if level <= max {
			out = append(out, p)
		}
	}
	return out
}

func parseVelocity(s string) VelocityLevel {
	switch s {
	case "slow":
		return VelocitySlow
	case "moderate":
		return VelocityModerate
	case "fast":
		return VelocityFast
	case "rapid":
		return VelocityRapid
	default:
		return VelocityNone
	}
}

func extractEventTimes(p scanner.PortInfo) []time.Time {
	if p.Meta == nil {
		return nil
	}
	raw, ok := p.Meta["event_times"]
	if !ok || raw == "" {
		return nil
	}
	var times []time.Time
	for _, s := range splitCSV(raw) {
		if ts, err := time.Parse(time.RFC3339, s); err == nil {
			times = append(times, ts)
		}
	}
	return times
}
