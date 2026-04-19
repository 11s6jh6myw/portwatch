package portmeta

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeVelocityPort(port int, eventTimes string) scanner.PortInfo {
	return scanner.PortInfo{
		Port: port,
		Meta: map[string]string{"event_times": eventTimes},
	}
}

func TestVelocityAnnotator_AddsMetadata(t *testing.T) {
	now := time.Now()
	evStr := now.Add(-2*time.Minute).Format(time.RFC3339)
	ports := []scanner.PortInfo{makeVelocityPort(80, evStr)}
	annotate := NewVelocityAnnotator(time.Hour)
	out := annotate(ports)
	if out[0].Meta[metaVelocity] == "" {
		t.Error("expected velocity metadata to be set")
	}
	if out[0].Meta[metaVelocityCount] == "" {
		t.Error("expected velocity_count metadata to be set")
	}
}

func TestVelocityAnnotator_NoEvents_NoneLevel(t *testing.T) {
	ports := []scanner.PortInfo{{Port: 443, Meta: map[string]string{}}}
	annotate := NewVelocityAnnotator(time.Hour)
	out := annotate(ports)
	if got := out[0].Meta[metaVelocity]; got != "none" {
		t.Errorf("expected none, got %q", got)
	}
}

func TestFilterByMaxVelocity_IncludesSlow(t *testing.T) {
	ports := []scanner.PortInfo{
		{Port: 80, Meta: map[string]string{metaVelocity: "slow"}},
		{Port: 443, Meta: map[string]string{metaVelocity: "rapid"}},
	}
	out := FilterByMaxVelocity(ports, VelocitySlow)
	if len(out) != 1 || out[0].Port != 80 {
		t.Errorf("expected only port 80, got %v", out)
	}
}

func TestFilterByMaxVelocity_NoMeta_PassesThrough(t *testing.T) {
	ports := []scanner.PortInfo{{Port: 22}}
	out := FilterByMaxVelocity(ports, VelocitySlow)
	if len(out) != 1 {
		t.Errorf("expected port with no meta to pass through")
	}
}
