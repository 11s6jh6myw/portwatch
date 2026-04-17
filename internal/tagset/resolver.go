package tagset

import (
	"os"
	"strings"
)

// Resolver enriches a Set with well-known automatic tags derived from
// the runtime environment.
type Resolver struct {
	hostname string
	env      map[string]string
}

// NewResolver creates a Resolver, reading the current hostname once.
func NewResolver() *Resolver {
	host, _ := os.Hostname()
	return &Resolver{
		hostname: host,
		env:      envMap(),
	}
}

// Enrich returns a new Set with host and any PORTWATCH_TAG_* env vars merged in.
// Caller-supplied tags take precedence over auto-detected ones.
func (r *Resolver) Enrich(s *Set) *Set {
	auto := New()
	if r.hostname != "" {
		auto.Add("host", r.hostname)
	}
	for k, v := range r.env {
		auto.Add(k, v)
	}
	// auto is the base; caller wins
	return auto.Merge(s)
}

func envMap() map[string]string {
	const prefix = "PORTWATCH_TAG_"
	m := make(map[string]string)
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, prefix) {
			continue
		}
		parts := strings.SplitN(e, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimPrefix(parts[0], prefix))
		if key != "" {
			m[key] = parts[1]
		}
	}
	return m
}
