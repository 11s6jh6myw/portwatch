package portmeta

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeNoisePort(port int) scanner.PortInfo {
	return scanner.PortInfo{Host: "localhost", Port: port}
}

func TestNoiseLevel_String(t *testing.T) {
	cases := []struct {
		lvl  NoiseLevel
		want string
	}{
		{NoiseNone, "none"},
		{NoiseMinimal, "minimal"},
		{NoiseModerate, "moderate"},
		{NoiseHigh, "high"},
	}
	for _, c := range cases {
		if got := c.lvl.String(); got != c.want {
			t.Errorf("NoiseLevel(%d).String() = %q, want %q", c.lvl, got, c.want)
		}
	}
}

func TestNoiseFor_NoEvents(t *testing.T) {
	p := makeNoisePort(80)
	if got := NoiseFor(p, nil, time.Hour); got != NoiseNone {
		t.Errorf("expected NoiseNone, got %s", got)
	}
}

func TestNoiseFor_LowCount(t *testing.T) {
	p := makeNoisePort(80)
	events := []time.Time{time.Now().Add(-1 * time.Minute)}
	if got := NoiseFor(p, events, time.Hour); got != NoiseMinimal {
		t.Errorf("expected NoiseMinimal, got %s", got)
	}
}

func TestNoiseFor_ModerateCount(t *testing.T) {
	p := makeNoisePort(80)
	events := make([]time.Time, 10)
	for i := range events {
		events[i] = time.Now().Add(-time.Duration(i) * time.Minute)
	}
	if got := NoiseFor(p, events, time.Hour); got != NoiseModerate {
		t.Errorf("expected NoiseModerate, got %s", got)
	}
}

func TestNoiseFor_HighCount(t *testing.T) {
	p := makeNoisePort(80)
	events := make([]time.Time, 25)
	for i := range events {
		events[i] = time.Now().Add(-time.Duration(i) * time.Minute)
	}
	if got := NoiseFor(p, events, time.Hour); got != NoiseHigh {
		t.Errorf("expected NoiseHigh, got %s", got)
	}
}

func TestNoiseAnnotator_AddsMetadata(t *testing.T) {
	p := makeNoisePort(443)
	events := map[int][]time.Time{
		443: {time.Now().Add(-5 * time.Minute), time.Now().Add(-10 * time.Minute)},
	}
	annotate := NewNoiseAnnotator(events, time.Hour)
	out := annotate([]scanner.PortInfo{p})
	if out[0].Metadata[noiseKey] == "" {
		t.Error("expected noise.level metadata to be set")
	}
	if out[0].Metadata[noiseCountKey] == "" {
		t.Error("expected noise.event_count metadata to be set")
	}
}

func TestFilterByMaxNoise_IncludesLow(t *testing.T) {
	p := makeNoisePort(22)
	p.Metadata = map[string]string{noiseKey: "minimal"}
	out := FilterByMaxNoise([]scanner.PortInfo{p}, NoiseModerate)
	if len(out) != 1 {
		t.Errorf("expected 1 port, got %d", len(out))
	}
}

func TestFilterByMaxNoise_ExcludesHigh(t *testing.T) {
	p := makeNoisePort(22)
	p.Metadata = map[string]string{noiseKey: "high"}
	out := FilterByMaxNoise([]scanner.PortInfo{p}, NoiseModerate)
	if len(out) != 0 {
		t.Errorf("expected 0 ports, got %d", len(out))
	}
}
