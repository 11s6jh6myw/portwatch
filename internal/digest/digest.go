// Package digest computes and compares fingerprints of port sets,
// allowing portwatch to detect whether the current scan differs from
// a previously recorded state without a full deep comparison.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// Digest is a hex-encoded SHA-256 fingerprint of a port set.
type Digest string

// Compute returns a stable Digest for the given ports.
// The digest is order-independent: ports are sorted before hashing.
func Compute(ports []scanner.PortInfo) Digest {
	keys := make([]string, 0, len(ports))
	for _, p := range ports {
		keys = append(keys, portKey(p))
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		_, _ = fmt.Fprintln(h, k)
	}
	return Digest(hex.EncodeToString(h.Sum(nil)))
}

// Equal reports whether two digests are identical.
func Equal(a, b Digest) bool {
	return a == b
}

// String returns the underlying hex string of the digest.
func (d Digest) String() string {
	return string(d)
}

// Short returns the first 8 characters of the digest, useful for
// compact display in logs or CLI output.
func (d Digest) Short() string {
	if len(d) <= 8 {
		return string(d)
	}
	return string(d[:8])
}

func portKey(p scanner.PortInfo) string {
	return fmt.Sprintf("%s:%d", p.Protocol, p.Port)
}
