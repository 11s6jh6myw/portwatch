// Package rollup aggregates bursts of port-change events into a single
// summary emitted after a configurable time window.
//
// Instead of firing one alert per port change during a large deployment
// or network reconfiguration, Roller collects all events that arrive
// within the window and delivers them together via an onFlush callback.
//
// Basic usage:
//
//	r := rollup.New(5*time.Second, func(s rollup.Summary) {
//		text, _ := rollup.Format(s, "text")
//		fmt.Print(text)
//	})
//	r.Add(rollup.Event{Port: p, Opened: true})
package rollup
