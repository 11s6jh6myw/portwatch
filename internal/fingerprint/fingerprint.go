// Package fingerprint produces a stable identity string for a set of open ports,
// allowing quick detection of scan-to-scan changes without a full diff.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// Fingerprint is a hex-encoded SHA-256 digest of a sorted port set.
type Fingerprint string

// Compute returns a Fingerprint for the given ports.
// The result is order-independent: the same set of ports always produces the
// same fingerprint regardless of the slice order.
func Compute(ports []scanner.PortInfo) Fingerprint {
	keys := make([]string, 0, len(ports))
	for _, p := range ports {
		keys = append(keys, portKey(p))
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintln(h, k)
	}
	return Fingerprint(hex.EncodeToString(h.Sum(nil)))
}

// Equal reports whether two fingerprints are identical.
func Equal(a, b Fingerprint) bool { return a == b }

func portKey(p scanner.PortInfo) string {
	return fmt.Sprintf("%s:%d:%s", p.Protocol, p.Port, p.State)
}
