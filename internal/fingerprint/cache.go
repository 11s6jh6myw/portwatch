package fingerprint

import "sync"

// Cache stores the most recent fingerprint per host so the monitor can skip
// a full diff when the fingerprint is unchanged.
type Cache struct {
	mu    sync.Mutex
	store map[string]Fingerprint
}

// NewCache returns an initialised Cache.
func NewCache() *Cache {
	return &Cache{store: make(map[string]Fingerprint)}
}

// Changed reports whether fp differs from the cached value for key.
// It always updates the cache to fp before returning.
func (c *Cache) Changed(key string, fp Fingerprint) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	prev, ok := c.store[key]
	c.store[key] = fp
	if !ok {
		return true
	}
	return !Equal(prev, fp)
}

// Get returns the cached fingerprint for key and whether it exists.
func (c *Cache) Get(key string) (Fingerprint, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f, ok := c.store[key]
	return f, ok
}

// Invalidate removes the cached entry for key.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}
