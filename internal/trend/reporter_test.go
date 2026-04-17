package trend_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/trend"
)

func makeTrend(port uint16, opened, closed int, flapping bool) trend.PortTrend {
	return trend.PortTrend{
		Port:       port,
		OpenCount:  opened,
		CloseCount: closed,
		Flapping:   flapping,
		FirstSeen:  base,
		LastSeen:   base.Add(time.Minute),
	}
}

func TestReport_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	trends := []trend.PortTrend{makeTrend(80, 3, 2, true)}
	if err := trend.Report(&buf, trends, "text"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "80") {
		t.Error("expected port 80 in output")
	}
	if !strings.Contains(out, "YES") {
		t.Error("expected flapping=YES in output")
	}
}

func TestReport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	trends := []trend.PortTrend{makeTrend(443, 1, 1, false)}
	if err := trend.Report(&buf, trends, "json"); err != nil {
		t.Fatal(err)
	}
	var out []trend.PortTrend
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(out) != 1 || out[0].Port != 443 {
		t.Errorf("unexpected output: %+v", out)
	}
}

func TestReport_SortedByPort(t *testing.T) {
	var buf bytes.Buffer
	trends := []trend.PortTrend{makeTrend(9000, 1, 0, false), makeTrend(80, 1, 0, false)}
	if err := trend.Report(&buf, trends, "text"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	i80 := strings.Index(out, "80")
	i9000 := strings.Index(out, "9000")
	if i80 > i9000 {
		t.Error("expected port 80 before 9000 in output")
	}
}
