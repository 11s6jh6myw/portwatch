// Package filter provides port filtering capabilities for portwatch.
// It allows users to exclude specific ports or port ranges from monitoring.
package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// Rule represents a single filter rule, either an exact port or a range.
type Rule struct {
	Low  uint16
	High uint16
}

// Filter holds a set of rules used to exclude ports from monitoring.
type Filter struct {
	rules []Rule
}

// New creates a Filter from a slice of rule strings.
// Each string can be a single port (e.g. "22") or a range (e.g. "8000-9000").
func New(specs []string) (*Filter, error) {
	f := &Filter{}
	for _, spec := range specs {
		rule, err := parseRule(spec)
		if err != nil {
			return nil, fmt.Errorf("filter: invalid rule %q: %w", spec, err)
		}
		f.rules = append(f.rules, rule)
	}
	return f, nil
}

// Excluded reports whether the given port should be excluded from monitoring.
func (f *Filter) Excluded(port uint16) bool {
	for _, r := range f.rules {
		if port >= r.Low && port <= r.High {
			return true
		}
	}
	return false
}

// Rules returns the list of active filter rules.
func (f *Filter) Rules() []Rule {
	out := make([]Rule, len(f.rules))
	copy(out, f.rules)
	return out
}

func parseRule(spec string) (Rule, error) {
	parts := strings.SplitN(spec, "-", 2)
	switch len(parts) {
	case 1:
		p, err := parsePort(parts[0])
		if err != nil {
			return Rule{}, err
		}
		return Rule{Low: p, High: p}, nil
	case 2:
		lo, err := parsePort(parts[0])
		if err != nil {
			return Rule{}, err
		}
		hi, err := parsePort(parts[1])
		if err != nil {
			return Rule{}, err
		}
		if lo > hi {
			return Rule{}, fmt.Errorf("low %d exceeds high %d", lo, hi)
		}
		return Rule{Low: lo, High: hi}, nil
	}
	return Rule{}, fmt.Errorf("unexpected format")
}

func parsePort(s string) (uint16, error) {
	s = strings.TrimSpace(s)
	n, err := strconv.ParseUint(s, 10, 16)
	if err != nil || n == 0 {
		return 0, fmt.Errorf("invalid port %q", s)
	}
	return uint16(n), nil
}
