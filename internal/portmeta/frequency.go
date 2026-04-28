package portmeta

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// FrequencyLevel describes how often a port has been observed open across scans.
type FrequencyLevel int

const (
	FrequencyNone FrequencyLevel = iota
	FrequencyRare
	FrequencyOccasional
	FrequencyCommon
	FrequencyAlways
)

func (f FrequencyLevel) String() string {
	switch f {
	case FrequencyRare:
		return "rare"
	case FrequencyOccasional:
		return "occasional"
	case FrequencyCommon:
		return "common"
	case FrequencyAlways:
		return "always"
	default:
		return "none"
	}
}

// FrequencyFor computes how frequently a port has been seen open relative to
// the total number of scans recorded in its metadata. It uses the meta keys
// "seen_count" and "scan_count" written by the lifecycle annotator.
func FrequencyFor(p scanner.PortInfo) FrequencyLevel {
	if p.Meta == nil {
		return FrequencyNone
	}

	seenRaw, ok := p.Meta["seen_count"]
	if !ok {
		return FrequencyNone
	}
	scanRaw, ok := p.Meta["scan_count"]
	if !ok {
		return FrequencyNone
	}

	seen := parseMetaInt(seenRaw)
	scans := parseMetaInt(scanRaw)
	if scans <= 0 || seen <= 0 {
		return FrequencyNone
	}

	ratio := float64(seen) / float64(scans)
	switch {
	case ratio >= 0.95:
		return FrequencyAlways
	case ratio >= 0.65:
		return FrequencyCommon
	case ratio >= 0.30:
		return FrequencyOccasional
	default:
		return FrequencyRare
	}
}

// IsFrequent returns true when the port is seen at least as often as
// FrequencyCommon.
func IsFrequent(p scanner.PortInfo) bool {
	return FrequencyFor(p) >= FrequencyCommon
}

// frequencyKey is used internally to avoid string literals across files.
const frequencyKey = "frequency"

// ensure time import is used if needed in future helpers.
var _ = time.Now
