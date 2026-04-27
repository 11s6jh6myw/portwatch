package portmeta

import (
	"strconv"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeRemPort(port int, firstSeen int64, events, scans int) scanner.PortInfo {
	return scanner.PortInfo{
		Port: port,
		Meta: map[string]string{
			"first_seen":  strconv.FormatInt(firstSeen, 10),
			"event_count": strconv.Itoa(events),
			"scan_count":  strconv.Itoa(scans),
		},
	}
}

func TestRemediationAnnotator_AddsMetadata(t *testing.T) {
	annotate := NewRemediationAnnotator()
	firstSeen := time.Now().Add(-10 * time.Minute).Unix()
	ports := []scanner.PortInfo{makeRemPort(23, firstSeen, 30, 50)}
	result := annotate(ports)

	if _, ok := result[0].Meta["remediation"]; !ok {
		t.Fatal("expected remediation key in meta")
	}
	if _, ok := result[0].Meta["remediation_actionable"]; !ok {
		t.Fatal("expected remediation_actionable key in meta")
	}
}

func TestRemediationAnnotator_HighRiskPort(t *testing.T) {
	annotate := NewRemediationAnnotator()
	firstSeen := time.Now().Add(-30 * time.Minute).Unix()
	ports := []scanner.PortInfo{makeRemPort(23, firstSeen, 40, 80)}
	result := annotate(ports)

	level := parseRemediation(result[0].Meta["remediation"])
	if level < RemediationReview {
		t.Errorf("expected >= review for high-risk port, got %s", level)
	}
}

func TestFilterByMinRemediation_IncludesActionable(t *testing.T) {
	annotate := NewRemediationAnnotator()
	firstSeen := time.Now().Add(-1 * time.Hour).Unix()
	ports := []scanner.PortInfo{
		makeRemPort(23, firstSeen, 50, 100),
		makeRemPort(443, time.Now().Add(-30*24*time.Hour).Unix(), 2, 200),
	}
	annotated := annotate(ports)
	filtered := FilterByMinRemediation(annotated, RemediationReview)

	for _, p := range filtered {
		level := parseRemediation(p.Meta["remediation"])
		if level < RemediationReview {
			t.Errorf("port %d should not pass filter, level=%s", p.Port, level)
		}
	}
}

func TestFilterByMinRemediation_NoMeta_PassesThrough(t *testing.T) {
	ports := []scanner.PortInfo{{Port: 8080}}
	result := FilterByMinRemediation(ports, RemediationNone)
	if len(result) != 1 {
		t.Errorf("expected 1 port, got %d", len(result))
	}
}
