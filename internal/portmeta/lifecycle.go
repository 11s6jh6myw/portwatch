package portmeta

import "time"

// LifecycleStage describes the observed stability of a port over time.
type LifecycleStage string

const (
	StageNew      LifecycleStage = "new"
	StageStable   LifecycleStage = "stable"
	StageFlapping LifecycleStage = "flapping"
	StageRetired  LifecycleStage = "retired"
)

func (s LifecycleStage) String() string { return string(s) }

// LifecycleEntry records the open/close history used to classify a port.
type LifecycleEntry struct {
	Port      int
	Protocol  string
	FirstSeen time.Time
	LastSeen  time.Time
	OpenCount int
	CloseCount int
}

// Classify returns a LifecycleStage based on observed open/close counts
// and how recently the port was last seen.
func Classify(e LifecycleEntry, now time.Time, stableAfter time.Duration) LifecycleStage {
	age := now.Sub(e.FirstSeen)

	switch {
	case e.CloseCount == 0 && e.OpenCount == 1:
		if age < stableAfter {
			return StageNew
		}
		return StageStable
	case e.CloseCount > 0 && now.Sub(e.LastSeen) > stableAfter && e.OpenCount == 0:
		return StageRetired
	case e.CloseCount >= 2 || (e.OpenCount > 1 && e.CloseCount >= 1):
		return StageFlapping
	case age >= stableAfter:
		return StageStable
	default:
		return StageNew
	}
}
