package portmeta

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

const (
	metaLifespan  = "lifespan"
	metaFirstSeen = "first_seen"
)

// NewLifespanAnnotator returns an annotator that adds lifespan metadata
// to each PortInfo based on the first_seen meta field.
func NewLifespanAnnotator() func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			firstSeen := parseMetaTime(p.Meta[metaFirstSeen])
			level := LifespanFor(firstSeen)
			if p.Meta == nil {
				p.Meta = make(map[string]string)
			}
			p.Meta[metaLifespan] = level.String()
			out[i] = p
		}
		return out
	}
}

// FilterByMinLifespan returns only ports whose lifespan meets or exceeds min.
// Ports without lifespan metadata are passed through unchanged.
func FilterByMinLifespan(ports []scanner.PortInfo, min LifespanLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		raw, ok := p.Meta[metaLifespan]
		if !ok {
			out = append(out, p)
			continue
		}
		if parseLifespan(raw) >= min {
			out = append(out, p)
		}
	}
	return out
}

func parseLifespan(s string) LifespanLevel {
	switch s {
	case "ephemeral":
		return LifespanEphemeral
	case "short":
		return LifespanShort
	case "medium":
		return LifespanMedium
	case "long":
		return LifespanLong
	case "permanent":
		return LifespanPermanent
	default:
		return LifespanUnknown
	}
}

// firstSeenFromMeta parses the first_seen meta field into a time.Time.
func firstSeenFromMeta(meta map[string]string) time.Time {
	return parseMetaTime(meta[metaFirstSeen])
}
