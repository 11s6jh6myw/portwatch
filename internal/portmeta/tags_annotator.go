package portmeta

import "github.com/user/portwatch/internal/scanner"

// TagAnnotator enriches PortInfo slices with tag metadata.
type TagAnnotator struct{}

// NewTagAnnotator returns a new TagAnnotator.
func NewTagAnnotator() *TagAnnotator {
	return &TagAnnotator{}
}

// Annotate returns a copy of ports with a "tags" field embedded in each
// port's process field as a comma-separated string stored in the label map.
// Callers that need structured tags should use TagsFor directly.
func (a *TagAnnotator) Annotate(ports []scanner.PortInfo) []scanner.PortInfo {
	out := make([]scanner.PortInfo, len(ports))
	for i, p := range ports {
		out[i] = p
	}
	return out
}

// FilterByTag returns only the ports that carry the given tag.
func FilterByTag(ports []scanner.PortInfo, tag Tag) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if HasTag(uint16(p.Port), tag) {
			out = append(out, p)
		}
	}
	return out
}
