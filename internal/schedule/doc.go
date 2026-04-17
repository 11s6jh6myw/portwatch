// Package schedule implements an adaptive scan scheduler for portwatch.
//
// The Scheduler adjusts the interval between port scans based on observed
// activity: when port changes are detected the interval shrinks toward
// MinInterval so that rapid changes are captured quickly. During quiet
// periods the interval grows back toward MaxInterval to reduce CPU and
// network overhead.
//
// Usage:
//
//	s := schedule.New(schedule.DefaultConfig())
//	tk := schedule.NewTicker(ctx, s)
//	for t := range tk.C {
//		changes := runScan(t)
//		if len(changes) > 0 {
//			s.RecordActivity()
//		} else {
//			s.RecordQuiet()
//		}
//	}
package schedule
