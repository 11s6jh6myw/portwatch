// Package reporter formats and writes port change summaries
// to various output destinations (stdout, file, etc.).
package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Format controls how reports are rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds a snapshot of port changes at a point in time.
type Report struct {
	Timestamp time.Time
	Opened    []scanner.PortInfo
	Closed    []scanner.PortInfo
}

// Reporter writes Reports to an io.Writer.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a Reporter writing to out using the given format.
// If out is nil, os.Stdout is used.
func New(out io.Writer, format Format) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{out: out, format: format}
}

// Write renders and writes a Report to the underlying writer.
func (r *Reporter) Write(rep Report) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(rep)
	default:
		return r.writeText(rep)
	}
}

func (r *Reporter) writeText(rep Report) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] Port changes detected\n", rep.Timestamp.Format(time.RFC3339)))
	for _, p := range rep.Opened {
		sb.WriteString(fmt.Sprintf("  + OPENED  %s\n", p))
	}
	for _, p := range rep.Closed {
		sb.WriteString(fmt.Sprintf("  - CLOSED  %s\n", p))
	}
	_, err := fmt.Fprint(r.out, sb.String())
	return err
}

func (r *Reporter) writeJSON(rep Report) error {
	opened := formatPortList(rep.Opened)
	closed := formatPortList(rep.Closed)
	line := fmt.Sprintf(
		`{"timestamp":%q,"opened":[%s],"closed":[%s]}\n`,
		rep.Timestamp.Format(time.RFC3339),
		strings.Join(opened, ","),
		strings.Join(closed, ","),
	)
	_, err := fmt.Fprint(r.out, line)
	return err
}

func formatPortList(ports []scanner.PortInfo) []string {
	out := make([]string, len(ports))
	for i, p := range ports {
		out[i] = fmt.Sprintf(`{"port":%d,"proto":%q}`, p.Port, p.Proto)
	}
	return out
}
