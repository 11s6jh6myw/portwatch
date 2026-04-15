package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event describes a port change event.
type Event struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Port      scanner.PortInfo
}

// Notifier sends alerts for port change events.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// NotifyOpened emits an alert for a newly opened port.
func (n *Notifier) NotifyOpened(p scanner.PortInfo) {
	n.emit(Event{
		Timestamp: time.Now(),
		Level:     LevelAlert,
		Message:   "port opened",
		Port:      p,
	})
}

// NotifyClosed emits an alert for a newly closed port.
func (n *Notifier) NotifyClosed(p scanner.PortInfo) {
	n.emit(Event{
		Timestamp: time.Now(),
		Level:     LevelWarn,
		Message:   "port closed",
		Port:      p,
	})
}

func (n *Notifier) emit(e Event) {
	fmt.Fprintf(n.out, "[%s] %s %s — %s\n",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Port.String(),
		e.Message,
	)
}
