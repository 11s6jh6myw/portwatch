package portmeta

import (
	"testing"
	"time"

	"github.com/netwatch/portwatch/internal/scanner"
)

func makeFreqPort(port int) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Proto: "tcp", State: "open"}
}

func TestChangeFreqAnnotator_AddsMetadata(t *testing.T) {
	a := NewChangeFreqAnnotator(time.Hour, func(port int) int { return 8 })
	ports := []scanner.PortInfo{makeFreqPort(80)}
	out := a.Annotate(ports)

	if out[0].Metadata[changeFreqKey] != "frequent" {
		t.Errorf("expected frequent, got %s", out[0].Metadata[changeFreqKey])
	}
	if out[0].Metadata[changeFreqCountKey] != "8" {
		t.Errorf("unexpected count: %s", out[0].Metadata[changeFreqCountKey])
	}
}

func TestChangeFreqAnnotator_StablePort(t *testing.T) {
	a := NewChangeFreqAnnotator(24*time.Hour, func(port int) int { return 0 })
	out := a.Annotate([]scanner.PortInfo{makeFreqPort(443)})
	if out[0].Metadata[changeFreqKey] != "stable" {
		t.Errorf("expected stable, got %s", out[0].Metadata[changeFreqKey])
	}
}

func TestFilterByMaxChangeFreq_IncludesStable(t *testing.T) {
	a := NewChangeFreqAnnotator(time.Hour, func(port int) int {
		if port == 80 {
			return 0
		}
		return 30
	})
	ports := a.Annotate([]scanner.PortInfo{makeFreqPort(80), makeFreqPort(8080)})
	out := FilterByMaxChangeFreq(ports, ChangeFreqStable)
	if len(out) != 1 || out[0].Port != 80 {
		t.Errorf("expected only port 80, got %v", out)
	}
}

func TestFilterByMaxChangeFreq_NoMetadata_PassesThrough(t *testing.T) {
	ports := []scanner.PortInfo{makeFreqPort(22)}
	out := FilterByMaxChangeFreq(ports, ChangeFreqStable)
	if len(out) != 1 {
		t.Errorf("expected port without metadata to pass through")
	}
}

func TestFilterByMaxChangeFreq_AllowsUpToMax(t *testing.T) {
	a := NewChangeFreqAnnotator(time.Hour, func(port int) int { return 8 }) // frequent
	ports := a.Annotate([]scanner.PortInfo{makeFreqPort(9000)})
	out := FilterByMaxChangeFreq(ports, ChangeFreqFrequent)
	if len(out) != 1 {
		t.Errorf("expected port to pass filter at max=frequent")
	}
}
