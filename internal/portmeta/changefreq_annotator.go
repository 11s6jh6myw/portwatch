package portmeta

import (
	"strconv"
	"time"

	"github.com/netwatch/portwatch/internal/scanner"
)

const (
	changeFreqKey        = "change_freq"
	changeFreqWindowKey  = "change_freq_window_hours"
	changeFreqCountKey   = "change_freq_transitions"
)

// ChangeFreqAnnotator attaches a change-frequency label to each PortInfo.
type ChangeFreqAnnotator struct {
	transitionsFn func(port int) int
	window        time.Duration
}

// NewChangeFreqAnnotator returns an annotator that calls transitionsFn to
// retrieve the number of observed state transitions for a port within window.
func NewChangeFreqAnnotator(window time.Duration, transitionsFn func(port int) int) *ChangeFreqAnnotator {
	return &ChangeFreqAnnotator{transitionsFn: transitionsFn, window: window}
}

// Annotate adds change_freq, change_freq_window_hours, and
// change_freq_transitions metadata fields to each port.
func (a *ChangeFreqAnnotator) Annotate(ports []scanner.PortInfo) []scanner.PortInfo {
	out := make([]scanner.PortInfo, len(ports))
	for i, p := range ports {
		n := a.transitionsFn(p.Port)
		cf := ClassifyChangeFreq(n, a.window)
		if p.Metadata == nil {
			p.Metadata = make(map[string]string)
		}
		p.Metadata[changeFreqKey] = cf.String()
		p.Metadata[changeFreqWindowKey] = strconv.FormatFloat(a.window.Hours(), 'f', 2, 64)
		p.Metadata[changeFreqCountKey] = strconv.Itoa(n)
		out[i] = p
	}
	return out
}

// FilterByMaxChangeFreq returns only ports whose change frequency is at or
// below the supplied maximum.
func FilterByMaxChangeFreq(ports []scanner.PortInfo, max ChangeFreq) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		raw, ok := p.Metadata[changeFreqKey]
		if !ok {
			out = append(out, p)
			continue
		}
		for cf := ChangeFreqStable; cf <= ChangeFreqVolatile; cf++ {
			if cf.String() == raw && cf <= max {
				out = append(out, p)
				break
			}
		}
	}
	return out
}
