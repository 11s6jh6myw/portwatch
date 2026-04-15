package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(port int) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Proto: "tcp"}
}

func TestNotifier_NotifyOpened(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	n.NotifyOpened(makePort(8080))

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT level in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
	if !strings.Contains(out, "port opened") {
		t.Errorf("expected 'port opened' message in output, got: %s", out)
	}
}

func TestNotifier_NotifyClosed(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	n.NotifyClosed(makePort(443))

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN level in output, got: %s", out)
	}
	if !strings.Contains(out, "443") {
		t.Errorf("expected port 443 in output, got: %s", out)
	}
	if !strings.Contains(out, "port closed") {
		t.Errorf("expected 'port closed' message in output, got: %s", out)
	}
}

func TestNotifier_DefaultsToStdout(t *testing.T) {
	// Should not panic when w is nil (falls back to os.Stdout).
	n := alert.New(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}
