// Package digest provides a lightweight fingerprinting mechanism for
// sets of open ports. It computes a stable, order-independent SHA-256
// digest that can be compared across scan cycles to quickly determine
// whether the port landscape has changed, avoiding unnecessary diff
// computation when nothing has moved.
package digest
