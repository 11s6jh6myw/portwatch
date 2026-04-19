package portmeta

import (
	"strconv"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// AgeAnnotator attaches an age classification to each PortInfo based on the
// first-seen timestamp stored in its metadata map.
type AgeAnnotator struct {
	now func() time.Time
}

// NewAgeAnnotator returns an AgeAnnotator that uses the real clock.
func NewAgeAnnotator() *AgeAnnotator {
	return &AgeAnnotator{now: time.Now}
}

// Annotate adds "age_class" and "age_seconds" keys to each port's Meta map.
func (a *AgeAnnotator) Annotate(ports []scanner.PortInfo) []scanner.PortInfo {
	now := a.now()
	for i, p := range ports {
		if p.Meta == nil {
			p.Meta = map[string]string{}
		}
		var firstSeen time.Time
		if raw, ok := p.Meta["first_seen"]; ok {
			if ts, err := time.Parse(time.RFC3339, raw); err == nil {
				firstSeen = ts
			}
		}
		cls := ClassifyAge(firstSeen, now)
		p.Meta["age_class"] = cls.String()
		if !firstSeen.IsZero() {
			p.Meta["age_seconds"] = strconv.FormatInt(int64(now.Sub(firstSeen).Seconds()), 10)
		}
		ports[i] = p
	}
	return ports
}

// FilterByMaxAge returns only ports whose age class is at most maxAge.
func FilterByMaxAge(ports []scanner.PortInfo, maxAge AgeClass) []scanner.PortInfo {
	out := ports[:0:0]
	for _, p := range ports {
		cls := AgeUnknown
		if p.Meta != nil {
			switch p.Meta["age_class"] {
			case "fresh":
				cls = AgeFresh
			case "short-lived":
				cls = AgeShortLived
			case "mature":
				cls = AgeMature
			case "established":
				cls = AgeEstablished
			}
		}
		if cls <= maxAge {
			out = append(out, p)
		}
	}
	return out
}
