// Package config provides loading and validation of portwatch configuration.
//
// Configuration is read from a YAML file. Unset fields fall back to defaults.
// Example configuration file:
//
//	# portwatch.yaml
//	scan_interval: 30s   # how often to scan ports (minimum 1s)
//	ports: [22, 80, 443] # ports to watch (empty = scan all)
//	alert_format: text   # output format: "text" or "json"
//	log_file: ""         # path to log file, empty = stderr
//
// Usage:
//
//	cfg, err := config.Load("portwatch.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
package config
