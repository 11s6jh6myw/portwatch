package labeler

import "github.com/user/portwatch/internal/scanner"

// EnrichedPort pairs a PortInfo with its resolved label.
type EnrichedPort struct {
	scanner.PortInfo
	Label string `json:"label"`
	Known bool   `json:"known"`
}

// Enricher wraps a Labeler and decorates slices of PortInfo.
type Enricher struct {
	l *Labeler
}

// NewEnricher returns an Enricher backed by the given Labeler.
func NewEnricher(l *Labeler) *Enricher {
	return &Enricher{l: l}
}

// Enrich returns a new slice of EnrichedPort for each entry in ports.
func (e *Enricher) Enrich(ports []scanner.PortInfo) []EnrichedPort {
	out := make([]EnrichedPort, len(ports))
	for i, p := range ports {
		out[i] = EnrichedPort{
			PortInfo: p,
			Label:    e.l.Label(uint16(p.Port)),
			Known:    e.l.Known(uint16(p.Port)),
		}
	}
	return out
}
