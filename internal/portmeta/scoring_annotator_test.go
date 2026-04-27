package portmeta

import (
	"strconv"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makeScoringPort(port int) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Proto: "tcp"}
}

func TestScoringAnnotator_AddsMetadata(t *testing.T) {
	annotate := NewScoringAnnotator()
	ports := []scanner.PortInfo{makeScoringPort(22), makeScoringPort(80)}
	result := annotate(ports)

	if len(result) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(result))
	}
	for _, p := range result {
		if p.Meta[metaScoreRaw] == "" {
			t.Errorf("port %d missing %s", p.Port, metaScoreRaw)
		}
		if p.Meta[metaScoreLevel] == "" {
			t.Errorf("port %d missing %s", p.Port, metaScoreLevel)
		}
		raw, err := strconv.Atoi(p.Meta[metaScoreRaw])
		if err != nil {
			t.Errorf("port %d: score_raw not numeric: %v", p.Port, err)
		}
		if raw < 0 || raw > 100 {
			t.Errorf("port %d: score_raw %d out of range", p.Port, raw)
		}
	}
}

func TestScoringAnnotator_HighRiskPort(t *testing.T) {
	annotate := NewScoringAnnotator()
	result := annotate([]scanner.PortInfo{makeScoringPort(23)})
	level := result[0].Meta[metaScoreLevel]
	if level == "negligible" || level == "low" {
		t.Errorf("telnet (23) expected high/critical score, got %q", level)
	}
}

func TestFilterByMinScore_IncludesHighEnough(t *testing.T) {
	ports := []scanner.PortInfo{
		makeScoringPort(23),  // high risk
		makeScoringPort(80),  // lower risk
	}
	high := FilterByMinScore(ports, ScoreHigh)
	for _, p := range high {
		if BucketScore(CompositeScoreFor(p)) < ScoreHigh {
			t.Errorf("port %d should not pass ScoreHigh filter", p.Port)
		}
	}
}

func TestFilterByMinScore_NoMeta_PassesThrough(t *testing.T) {
	ports := []scanner.PortInfo{makeScoringPort(80)}
	result := FilterByMinScore(ports, ScoreNegligible)
	if len(result) != 1 {
		t.Errorf("expected 1 port, got %d", len(result))
	}
}
