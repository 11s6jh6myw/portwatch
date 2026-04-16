package eventlog_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/eventlog"
	"github.com/user/portwatch/internal/scanner"
)

func TestQuery_FilterByKind(t *testing.T) {
	s := eventlog.NewStore(tempPath(t))
	p := samplePort()
	_ = s.Append("opened", p)
	_ = s.Append("closed", p)

	events, err := s.Query(eventlog.Filter{Kind: "opened"})
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Kind != "opened" {
		t.Errorf("expected 1 opened event, got %v", events)
	}
}

func TestQuery_FilterByPort(t *testing.T) {
	s := eventlog.NewStore(tempPath(t))
	_ = s.Append("opened", scanner.PortInfo{Port: 80, Proto: "tcp", State: "open"})
	_ = s.Append("opened", scanner.PortInfo{Port: 443, Proto: "tcp", State: "open"})

	events, err := s.Query(eventlog.Filter{Port: 443})
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Port.Port != 443 {
		t.Errorf("expected port 443, got %v", events)
	}
}

func TestQuery_FilterBySince(t *testing.T) {
	s := eventlog.NewStore(tempPath(t))
	_ = s.Append("opened", samplePort())
	cutoff := time.Now().UTC()
	time.Sleep(2 * time.Millisecond)
	_ = s.Append("closed", samplePort())

	events, err := s.Query(eventlog.Filter{Since: cutoff})
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Kind != "closed" {
		t.Errorf("expected 1 recent event, got %v", events)
	}
}

func TestQuery_NoFilter_ReturnsAll(t *testing.T) {
	s := eventlog.NewStore(tempPath(t))
	_ = s.Append("opened", samplePort())
	_ = s.Append("closed", samplePort())

	events, err := s.Query(eventlog.Filter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}
