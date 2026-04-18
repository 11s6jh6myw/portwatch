package portmeta

import "github.com/user/portwatch/internal/scanner"

// ExposureAnnotator attaches an ExposureLevel tag to each PortInfo.
type ExposureAnnotator struct{}

// NewExposureAnnotator returns a new ExposureAnnotator.
func NewExposureAnnotator() *ExposureAnnotator {
	return &ExposureAnnotator{}
}

// Annotate returns a copy of ports with an "exposure" field set.
func (a *ExposureAnnotator) Annotate(ports []scanner.PortInfo) []AnnotatedPort {
	out := make([]AnnotatedPort, len(ports))
	for i, p := range ports {
		out[i] = AnnotatedPort{
			PortInfo: p,
			Exposure: ExposureFor(p.Port),
		}
	}
	return out
}

// AnnotatedPort pairs a PortInfo with its computed ExposureLevel.
type AnnotatedPort struct {
	scanner.PortInfo
	Exposure ExposureLevel
}

// FilterByMinExposure returns only those AnnotatedPorts whose exposure is
// at or above the given minimum level.
func FilterByMinExposure(ports []AnnotatedPort, min ExposureLevel) []AnnotatedPort {
	var out []AnnotatedPort
	for _, p := range ports {
		if p.Exposure >= min {
			out = append(out, p)
		}
	}
	return out
}
