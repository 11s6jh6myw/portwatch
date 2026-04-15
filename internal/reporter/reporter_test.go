package reporter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/scanner"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func makeReport(opened, closed []scanner.PortInfo) reporter.Report {
	return reporter.Report{
		Timestamp: fixedTime,
		Opened:    opened,
		Closed:    closed,
	}
}

func TestReporter_TextFormat_Opened(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	rep := makeReport([]scanner.PortInfo{{Port: 8080, Proto: "tcp"}}, nil)

	if err := r.Write(rep); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ OPENED") {
		t.Errorf("expected OPENED marker, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
}

func TestReporter_TextFormat_Closed(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	rep := makeReport(nil, []scanner.PortInfo{{Port: 22, Proto: "tcp"}})

	if err := r.Write(rep); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "- CLOSED") {
		t.Errorf("expected CLOSED marker, got: %s", out)
	}
}

func TestReporter_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)
	rep := makeReport(
		[]scanner.PortInfo{{Port: 443, Proto: "tcp"}},
		[]scanner.PortInfo{{Port: 80, Proto: "tcp"}},
	)

	if err := r.Write(rep); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"port":443`) {
		t.Errorf("expected port 443 in JSON, got: %s", out)
	}
	if !strings.Contains(out, `"port":80`) {
		t.Errorf("expected port 80 in JSON, got: %s", out)
	}
	if !strings.Contains(out, `"timestamp"`) {
		t.Errorf("expected timestamp in JSON, got: %s", out)
	}
}

func TestReporter_DefaultsToStdout(t *testing.T) {
	// Ensure New does not panic with nil writer
	r := reporter.New(nil, reporter.FormatText)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
