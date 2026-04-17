package tagset

import (
	"fmt"
	"sort"
	"strings"
)

// Tag is a key=value label attached to a port event.
type Tag struct {
	Key   string
	Value string
}

func (t Tag) String() string { return fmt.Sprintf("%s=%s", t.Key, t.Value) }

// Set holds an ordered, deduplicated collection of Tags.
type Set struct {
	tags map[string]string
}

// New creates a Set from zero or more "key=value" strings.
// Malformed entries are silently ignored.
func New(raw ...string) *Set {
	s := &Set{tags: make(map[string]string)}
	for _, r := range raw {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) == 2 && parts[0] != "" {
			s.tags[parts[0]] = parts[1]
		}
	}
	return s
}

// Add inserts or overwrites a tag.
func (s *Set) Add(key, value string) {
	if key != "" {
		s.tags[key] = value
	}
}

// Get returns the value for key and whether it was present.
func (s *Set) Get(key string) (string, bool) {
	v, ok := s.tags[key]
	return v, ok
}

// All returns all tags sorted by key.
func (s *Set) All() []Tag {
	keys := make([]string, 0, len(s.tags))
	for k := range s.tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]Tag, len(keys))
	for i, k := range keys {
		out[i] = Tag{Key: k, Value: s.tags[k]}
	}
	return out
}

// String renders the set as a comma-separated list.
func (s *Set) String() string {
	tags := s.All()
	parts := make([]string, len(tags))
	for i, t := range tags {
		parts[i] = t.String()
	}
	return strings.Join(parts, ",")
}

// Merge returns a new Set combining s and other; other wins on conflict.
func (s *Set) Merge(other *Set) *Set {
	out := New()
	for k, v := range s.tags {
		out.tags[k] = v
	}
	for k, v := range other.tags {
		out.tags[k] = v
	}
	return out
}
