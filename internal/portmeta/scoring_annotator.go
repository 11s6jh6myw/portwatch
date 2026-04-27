package portmeta

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

const (
	metaScoreRaw   = "score_raw"
	metaScoreLevel = "score_level"
)

// NewScoringAnnotator returns an annotator that adds composite score metadata
// to each PortInfo under the keys "score_raw" and "score_level".
func NewScoringAnnotator() func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			raw := CompositeScoreFor(p)
			level := BucketScore(raw)
			if p.Meta == nil {
				p.Meta = make(map[string]string)
			}
			p.Meta[metaScoreRaw] = fmt.Sprintf("%d", raw)
			p.Meta[metaScoreLevel] = level.String()
			out[i] = p
		}
		return out
	}
}

// FilterByMinScore returns only ports whose composite score meets or exceeds
// the given ScoreLevel.
func FilterByMinScore(ports []scanner.PortInfo, min ScoreLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if BucketScore(CompositeScoreFor(p)) >= min {
			out = append(out, p)
		}
	}
	return out
}
