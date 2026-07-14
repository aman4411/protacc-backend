// Package cache provides a small, in-process, TTL-based cache with single-flight
// loading. It is designed for a single-instance deployment (e.g. Render free tier)
// where an in-memory cache removes almost all repeated read load from the database
// without any external infrastructure.
//
// Typical use is cache-aside via the generic Load helper:
//
//	svcs, err := cache.Load(c, "services:list", 10*time.Minute, func() ([]Service, error) {
//	    return repo.GetServices(ctx)
//	})
//
// Concurrent misses for the same key are collapsed into a single loader call
// (single-flight), which prevents a cache-miss stampede from hammering the DB
// under high traffic.
package cache

import (
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

const (
	// defaultMaxEntries bounds memory use. The set of distinct keys for this app
	// is small (catalog data keyed by slug/id), so this is rarely reached.
	defaultMaxEntries = 10000
	// janitorInterval is how often expired entries are swept out of the map.
	janitorInterval = time.Minute
)

type entry struct {
	value     any
	expiresAt time.Time
}

// Cache is a TTL, size-bounded, in-process cache safe for concurrent use.
type Cache struct {
	mu         sync.RWMutex
	store      map[string]entry
	maxEntries int
	group      singleflight.Group
}

// New creates a cache and starts a background janitor that purges expired entries.
func New() *Cache {
	c := &Cache{
		store:      make(map[string]entry),
		maxEntries: defaultMaxEntries,
	}
	go c.janitor()
	return c
}

func (c *Cache) janitor() {
	ticker := time.NewTicker(janitorInterval)
	defer ticker.Stop()
	for range ticker.C {
		c.purgeExpired()
	}
}

func (c *Cache) purgeExpired() {
	now := time.Now()
	c.mu.Lock()
	for k, e := range c.store {
		if now.After(e.expiresAt) {
			delete(c.store, k)
		}
	}
	c.mu.Unlock()
}

// Get returns the cached value for key if present and not expired.
func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	e, ok := c.store[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.value, true
}

// Set stores value under key for the given ttl.
func (c *Cache) Set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Cheap guard against unbounded growth: at capacity, drop expired entries
	// first, then evict arbitrary entries until back under the cap.
	if len(c.store) >= c.maxEntries {
		now := time.Now()
		for k, e := range c.store {
			if now.After(e.expiresAt) {
				delete(c.store, k)
			}
		}
		for k := range c.store {
			if len(c.store) < c.maxEntries {
				break
			}
			delete(c.store, k)
		}
	}
	c.store[key] = entry{value: value, expiresAt: time.Now().Add(ttl)}
}

// Invalidate removes a single key.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	delete(c.store, key)
	c.mu.Unlock()
}

// InvalidatePrefix removes every key that begins with prefix. Used to bust a
// whole family of cached reads (e.g. all "services:*") after an admin write.
func (c *Cache) InvalidatePrefix(prefix string) {
	c.mu.Lock()
	for k := range c.store {
		if strings.HasPrefix(k, prefix) {
			delete(c.store, k)
		}
	}
	c.mu.Unlock()
}

// Load returns the cached value for key, or computes it via loader (guarded by
// single-flight so concurrent misses trigger only one load), caches it, and
// returns it. A nil cache bypasses caching entirely, so callers can be wired
// with or without a cache. Loader errors are not cached.
func Load[T any](c *Cache, key string, ttl time.Duration, loader func() (T, error)) (T, error) {
	if c == nil {
		return loader()
	}
	if v, ok := c.Get(key); ok {
		if tv, ok := v.(T); ok {
			return tv, nil
		}
	}
	res, err, _ := c.group.Do(key, func() (any, error) {
		// Double-check: another goroutine may have populated it while we waited.
		if v, ok := c.Get(key); ok {
			if tv, ok := v.(T); ok {
				return tv, nil
			}
		}
		v, err := loader()
		if err != nil {
			return v, err
		}
		c.Set(key, v, ttl)
		return v, nil
	})
	if err != nil {
		var zero T
		return zero, err
	}
	return res.(T), nil
}
