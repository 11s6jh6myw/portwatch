package portmeta

import "github.com/example/portwatch/internal/scanner"

// CategoryAnnotator attaches a category tag to each PortInfo.
type CategoryAnnotator struct{}

// NewCategoryAnnotator returns a CategoryAnnotator.
func NewCategoryAnnotator() *CategoryAnnotator { return &CategoryAnnotator{} }

// Annotate returns a copy of ports with a "category" tag added.
func (a *CategoryAnnotator) Annotate(ports []scanner.PortInfo) []scanner.PortInfo {
	out := make([]scanner.PortInfo, len(ports))
	for i, p := range ports {
		cp := p
		cat := Categorize(uint16(cp.Port))
		if cp.Tags == nil {
			cp.Tags = map[string]string{}
		}
		cp.Tags["category"] = cat.String()
		out[i] = cp
	}
	return out
}

// FilterByCategory returns only ports whose category matches one of the
// supplied categories. If no categories are given all ports are returned.
func FilterByCategory(ports []scanner.PortInfo, cats ...Category) []scanner.PortInfo {
	if len(cats) == 0 {
		return ports
	}
	set := make(map[Category]struct{}, len(cats))
	for _, c := range cats {
		set[c] = struct{}{}
	}
	var out []scanner.PortInfo
	for _, p := range ports {
		if _, ok := set[Categorize(uint16(p.Port))]; ok {
			out = append(out, p)
		}
	}
	return out
}
