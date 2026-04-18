package labeler_test

import (
	"testing"

	"github.com/user/portwatch/internal/labeler"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(port int) scanner.PortInfo {
	return scanner.PortInfo{Host: "127.0.0.1", Port: port, Proto: "tcp"}
}

func TestEnrich_AddsLabelAndKnown(t *testing.T) {
	l := labeler.New(nil)
	e := labeler.NewEnricher(l)

	ports := []scanner.PortInfo{makePort(22), makePort(9999)}
	result := e.Enrich(ports)

	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}

	if result[0].Label != "ssh" {
		t.Errorf("expected ssh, got %s", result[0].Label)
	}
	if !result[0].Known {
		t.Error("expected port 22 to be known")
	}
	if result[1].Label != "9999" {
		t.Errorf("expected 9999, got %s", result[1].Label)
	}
	if result[1].Known {
		t.Error("expected port 9999 to be unknown")
	}
}

func TestEnrich_EmptyInput(t *testing.T) {
	l := labeler.New(nil)
	e := labeler.NewEnricher(l)
	result := e.Enrich(nil)
	if len(result) != 0 {
		t.Fatalf("expected empty result")
	}
}

func TestEnrich_OverrideApplied(t *testing.T) {
	l := labeler.New(map[uint16]string{8080: "gateway"})
	e := labeler.NewEnricher(l)
	result := e.Enrich([]scanner.PortInfo{makePort(8080)})
	if result[0].Label != "gateway" {
		t.Errorf("expected gateway, got %s", result[0].Label)
	}
}
