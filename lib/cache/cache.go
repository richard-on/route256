// Package cache provides a wrapper for a cache implementation.
package cache

import "gitlab.ozon.dev/rragusskiy/homework-1/lib/cache/metrics"

const (
	Hit  = "cache_hit"
	Miss = "cache_miss"

	Set = "cache_set"
	Get = "cache_get"

	SetOK = "cache_set_ok"
)

// Cacher is an interface for a cache implementation.
type Cacher[K comparable, V any] interface {
	Get(key K) (V, bool)
	Set(key K, value V)
}

// Cache is a wrapper for a cache implementation.
type Cache[K comparable, V any] struct {
	service string
	cacher  Cacher[K, V]
}

// NewCache creates a new cache wrapper.
func NewCache[K comparable, V any](service string, cacher Cacher[K, V]) *Cache[K, V] {
	return &Cache[K, V]{
		service: service,
		cacher:  cacher,
	}
}

// Get returns a value from the cache and collects prometheus metrics.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	metrics.RequestCounter.WithLabelValues(c.service, Get).Inc()

	resp, ok := c.cacher.Get(key)
	if !ok {
		metrics.HistogramResponseTime.WithLabelValues(c.service, Get, Miss)
		metrics.ResponseCounter.WithLabelValues(c.service, Get, Miss)
		return resp, false
	}

	metrics.HistogramResponseTime.WithLabelValues(c.service, Get, Hit)
	metrics.ResponseCounter.WithLabelValues(c.service, Get, Hit)
	return resp, true
}

// Set sets a value to the cache and collects prometheus metrics.
func (c *Cache[K, V]) Set(key K, value V) {
	metrics.RequestCounter.WithLabelValues(c.service, Set).Inc()

	c.cacher.Set(key, value)

	metrics.HistogramResponseTime.WithLabelValues(c.service, Set, SetOK)
}
