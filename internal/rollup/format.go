package rollup

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Format renders a Summary as either "text" or "json".
func Format(s Summary, format string) (string, error) {
	switch strings.ToLower(format) {
	case "json":
		return formatJSON(s)
	default:
		return formatText(s), nil
	}
}

func formatText(s Summary) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("[rollup] %s\n", s.At.Format("2006-01-02 15:04:05")))
	if len(s.Opened) > 0 {
		b.WriteString(fmt.Sprintf("  opened (%d):", len(s.Opened)))
		for _, p := range s.Opened {
			b.WriteString(fmt.Sprintf(" %d/%s", p.Port, p.Proto))
		}
		b.WriteString("\n")
	}
	if len(s.Closed) > 0 {
		b.WriteString(fmt.Sprintf("  closed (%d):", len(s.Closed)))
		for _, p := range s.Closed {
			b.WriteString(fmt.Sprintf(" %d/%s", p.Port, p.Proto))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func formatJSON(s Summary) (string, error) {
	type payload struct {
		At     string `json:"at"`
		Opened []int  `json:"opened"`
		Closed []int  `json:"closed"`
	}
	p := payload{At: s.At.Format("2006-01-02T15:04:05Z07:00")}
	for _, port := range s.Opened {
		p.Opened = append(p.Opened, port.Port)
	}
	for _, port := range s.Closed {
		p.Closed = append(p.Closed, port.Port)
	}
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
