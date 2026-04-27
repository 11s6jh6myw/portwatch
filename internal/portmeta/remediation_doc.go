// Package portmeta — remediation
//
// RemediationFor computes the recommended remediation level for a port based
// on a composite view of its risk score, anomaly level, and severity. The
// result is one of:
//
//   - none      – no action required
//   - monitor   – keep an eye on the port
//   - review    – investigate during next maintenance window
//   - mitigate  – schedule prompt remediation
//   - immediate – requires urgent human intervention
//
// NewRemediationAnnotator attaches "remediation" and "remediation_actionable"
// keys to each PortInfo.Meta map. FilterByMinRemediation can then be used to
// narrow a port list to only those that meet or exceed a given threshold.
package portmeta
