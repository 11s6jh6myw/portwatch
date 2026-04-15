// Package baseline provides functionality for saving and comparing a trusted
// snapshot of open ports (the "baseline") against live scan results.
//
// Workflow:
//
//  1. On first run (or when the user requests a reset), call Store.Save to
//     persist the current set of open ports as the trusted baseline.
//
//  2. On subsequent runs, call Store.Load to retrieve the saved baseline, then
//     pass it together with the latest scan results to Check.
//
//  3. Check returns a slice of Violation values describing ports that are
//     unexpectedly open or unexpectedly closed relative to the baseline.
//
// The baseline file is stored as JSON and is human-readable, making it easy
// to audit or manually adjust outside of portwatch.
package baseline
