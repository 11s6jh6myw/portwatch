package portmeta

import (
	"strconv"

	"github.com/user/portwatch/internal/scanner"
)

// NewRemediationAnnotator returns an annotator that attaches remediation
// metadata to each PortInfo entry.
func NewRemediationAnnotator() func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		for i, p := range ports {
			var firstSeen int64
			var eventCount, scanCount int

			if v, ok := p.Meta["first_seen"]; ok {
				firstSeen, _ = strconv.ParseInt(v, 10, 64)
			}
			if v, ok := p.Meta["event_count"]; ok {
				eventCount, _ = strconv.Atoi(v)
			}
			if v, ok := p.Meta["scan_count"]; ok {
				scanCount, _ = strconv.Atoi(v)
			}

			level := RemediationFor(p.Port, firstSeen, eventCount, scanCount)

			if p.Meta == nil {
				p.Meta = make(map[string]string)
			}
			p.Meta["remediation"] = level.String()
			p.Meta["remediation_actionable"] = strconv.FormatBool(IsActionable(level))
			ports[i] = p
		}
		return ports
	}
}

// FilterByMinRemediation returns only ports whose remediation level meets
// the minimum threshold.
func FilterByMinRemediation(ports []scanner.PortInfo, min RemediationLevel) []scanner.PortInfo {
	out := ports[:0:0]
	for _, p := range ports {
		level := parseRemediation(p.Meta["remediation"])
		if level >= min {
			out = append(out, p)
		}
	}
	return out
}

func parseRemediation(s string) RemediationLevel {
	switch s {
	case "monitor":
		return RemediationMonitor
	case "review":
		return RemediationReview
	case "mitigate":
		return RemediationMitigate
	case "immediate":
		return RemediationImmediate
	default:
		return RemediationNone
	}
}
