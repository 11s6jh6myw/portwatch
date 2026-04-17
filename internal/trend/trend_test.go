package trend_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/trend"
)

var base = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func TestAnalyze_NoEvents(t *testing.T) {
	a := trend.New(time.Hour, 3)
	results := a.Analyze(base)
	if len(results) != 0 {
		t.Fatalf("expected 0 trends, got %d", len(results))
	}
}

func TestAnalyze_CountsOpenAndClose(t *testing.T) {
	a := trend.New(time.Hour, 10)
	a.Record(trend.Event{Port: 80, Kind: "opened", Timestamp: base.Add(-10 * time.Minute)})
	a.Record(trend.Event{Port: 80, Kind: "closed", Timestamp: base.Add(-5 * time.Minute)})

	results := a.Analyze(base)
	if len(results) != 1 {
		t.Fatalf("expected 1 trend, got %d", len(results))
	}
	tr := results[0]
	if tr.OpenCount != 1 || tr.CloseCount != 1 {
		t.Errorf("unexpected counts: open=%d close=%d", tr.OpenCount, tr.CloseCount)
	}
}

func TestAnalyze_FlappingDetected(t *testing.T) {
	a := trend.New(time.Hour, 3)
	for i := 0; i < 2; i++ {
		a.Record(trend.Event{Port: 443, Kind: "opened", Timestamp: base.Add(-time.Duration(i+1) * time.Minute)})
		a.Record(trend.Event{Port: 443, Kind: "closed", Timestamp: base.Add(-time.Duration(i+1) * time.Minute)})
	}
	results := a.Analyze(base)
	if !results[0].Flapping {
		t.Error("expected port 443 to be marked as flapping")
	}
}

func TestAnalyze_EventsOutsideWindowIgnored(t *testing.T) {
	a := trend.New(time.Hour, 1)
	a.Record(trend.Event{Port: 22, Kind: "opened", Timestamp: base.Add(-2 * time.Hour)})
	results := a.Analyze(base)
	if len(results) != 0 {
		t.Errorf("expected 0 trends, got %d", len(results))
	}
}

func TestAnalyze_MultiplePorts(t *testing.T) {
	a := trend.New(time.Hour, 5)
	a.Record(trend.Event{Port: 80, Kind: "opened", Timestamp: base.Add(-1 * time.Minute)})
	a.Record(trend.Event{Port: 443, Kind: "opened", Timestamp: base.Add(-2 * time.Minute)})
	results := a.Analyze(base)
	if len(results) != 2 {
		t.Errorf("expected 2 trends, got %d", len(results))
	}
}
