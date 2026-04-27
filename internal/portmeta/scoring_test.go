package portmeta

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestBucketScore_Boundaries(t *testing.T) {
	cases := []struct {
		score int
		want  ScoreLevel
	}{
		{0, ScoreNegligible},
		{19, ScoreNegligible},
		{20, ScoreLow},
		{39, ScoreLow},
		{40, ScoreMedium},
		{59, ScoreMedium},
		{60, ScoreHigh},
		{79, ScoreHigh},
		{80, ScoreCritical},
		{100, ScoreCritical},
	}
	for _, tc := range cases {
		got := BucketScore(tc.score)
		if got != tc.want {
			t.Errorf("BucketScore(%d) = %v, want %v", tc.score, got, tc.want)
		}
	}
}

func TestScoreLevel_String(t *testing.T) {
	cases := []struct {
		level ScoreLevel
		want  string
	}{
		{ScoreNegligible, "negligible"},
		{ScoreLow, "low"},
		{ScoreMedium, "medium"},
		{ScoreHigh, "high"},
		{ScoreCritical, "critical"},
		{ScoreLevel(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("ScoreLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestNormalisedScore_ReturnsLevel(t *testing.T) {
	// Port 23 (telnet) should score high/critical.
	p := scanner.PortInfo{Port: 23, Proto: "tcp"}
	level := NormalisedScore(p)
	if level < ScoreMedium {
		t.Errorf("expected at least Medium for port 23, got %v", level)
	}
}
