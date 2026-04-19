package portmeta

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeChurnPort(port int) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Meta: make(map[string]string)}
}

func TestChurnAnnotator_AddsMetadata(t *testing.T) {
	events := map[int][]ChurnEvent{
		8080: {
			{At: time.Now().Add(-10 * time.Minute), Kind: "opened"},
			{At: time.Now().Add(-5 * time.Minute), Kind: "closed"},
		},
	}
	a := NewChurnAnnotator(events, time.Hour)
	ports := a.Annotate([]scanner.PortInfo{makeChurnPort(8080)})
	if ports[0].Meta[churnLevelKey] != "low" {
		t.Errorf("expected low churn, got %s", ports[0].Meta[churnLevelKey])
	}
	if ports[0].Meta[churnCountKey] != "2" {
		t.Errorf("expected count 2, got %s", ports[0].Meta[churnCountKey])
	}
}

func TestChurnAnnotator_NoEvents_NoneLevel(t *testing.T) {
	a := NewChurnAnnotator(map[int][]ChurnEvent{}, time.Hour)
	ports := a.Annotate([]scanner.PortInfo{makeChurnPort(443)})
	if ports[0].Meta[churnLevelKey] != "none" {
		t.Errorf("expected none, got %s", ports[0].Meta[churnLevelKey])
	}
}

func TestFilterByMaxChurn_IncludesLow(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 80, Meta: map[string]string{churnLevelKey: "low"}},
		{Port: 9000, Meta: map[string]string{churnLevelKey: "high"}},
	}
	got := FilterByMaxChurn(ports, ChurnLow)
	if len(got) != 1 || got[0].Port != 80 {
		t.Errorf("expected only port 80, got %+v", got)
	}
}

func TestFilterByMaxChurn_NoMeta_PassesThrough(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 22, Meta: nil},
	}
	got := FilterByMaxChurn(ports, ChurnNone)
	if len(got) != 1 {
		t.Errorf("expected port with no meta to pass through")
	}
}
