package portmeta

import "github.com/user/portwatch/internal/scanner"

// OwnerAnnotator attaches owner metadata to PortInfo values.
type OwnerAnnotator struct{}

// NewOwnerAnnotator returns a new OwnerAnnotator.
func NewOwnerAnnotator() *OwnerAnnotator { return &OwnerAnnotator{} }

// Annotate returns a copy of ports with owner information added to each
// port's metadata map under the keys "owner_org" and "owner_contact".
// Ports without a known owner are left unchanged.
func (a *OwnerAnnotator) Annotate(ports []scanner.PortInfo) []scanner.PortInfo {
	out := make([]scanner.PortInfo, len(ports))
	for i, p := range ports {
		if o, ok := LookupOwner(uint16(p.Port)); ok {
			if p.Meta == nil {
				p.Meta = make(map[string]string)
			}
			p.Meta["owner_org"] = o.Org
			p.Meta["owner_contact"] = o.Contact
		}
		out[i] = p
	}
	return out
}

// FilterByKnownOwner returns only ports that have a registered owner.
func FilterByKnownOwner(ports []scanner.PortInfo) []scanner.PortInfo {
	var out []scanner.PortInfo
	for _, p := range ports {
		if KnownOwner(uint16(p.Port)) {
			out = append(out, p)
		}
	}
	return out
}
