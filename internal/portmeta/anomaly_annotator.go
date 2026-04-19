package portmeta

import "github.com/iamcathal/portwatch/internal/scanner"

// NewAnomalyAnnotator returns an annotator that adds anomaly metadata to each port.
func NewAnomalyAnnotator() func([]scanner.PortInfo) []scanner.PortInfo {
	return func(ports []scanner.PortInfo) []scanner.PortInfo {
		out := make([]scanner.PortInfo, len(ports))
		for i, p := range ports {
			level := AnomalyFor(p.Port)
			if p.Metadata == nil {
				p.Metadata = make(map[string]string)
			}
			p.Metadata["anomaly"] = level.String()
			out[i] = p
		}
		return out
	}
}

// FilterByMinAnomaly returns only ports whose anomaly level is >= min.
func FilterByMinAnomaly(ports []scanner.PortInfo, min AnomalyLevel) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if AnomalyFor(p.Port) >= min {
			out = append(out, p)
		}
	}
	return out
}
