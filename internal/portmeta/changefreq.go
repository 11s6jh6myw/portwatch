package portmeta

import "time"

// ChangeFreq describes how often a port's state changes.
type ChangeFreq int

const (
	ChangeFreqStable   ChangeFreq = iota // rarely or never changes
	ChangeFreqOccasional                 // changes a few times per day
	ChangeFreqFrequent                   // changes many times per day
	ChangeFreqVolatile                   // changes continuously
)

func (c ChangeFreq) String() string {
	switch c {
	case ChangeFreqStable:
		return "stable"
	case ChangeFreqOccasional:
		return "occasional"
	case ChangeFreqFrequent:
		return "frequent"
	case ChangeFreqVolatile:
		return "volatile"
	default:
		return "unknown"
	}
}

// ClassifyChangeFreq returns a ChangeFreq based on the number of state
// transitions observed within the supplied observation window.
func ClassifyChangeFreq(transitions int, window time.Duration) ChangeFreq {
	if window <= 0 || transitions == 0 {
		return ChangeFreqStable
	}
	perHour := float64(transitions) / window.Hours()
	switch {
	case perHour >= 20:
		return ChangeFreqVolatile
	case perHour >= 5:
		return ChangeFreqFrequent
	case perHour >= 1:
		return ChangeFreqOccasional
	default:
		return ChangeFreqStable
	}
}
