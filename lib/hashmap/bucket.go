package hashmap

import "sync"

// bucket is a wrapper around a stdlib map.
type bucket[K comparable, V any] struct {
	// hmap is the underlying map.
	hmap map[K]V

	mu sync.RWMutex
}

// newBucket returns a new bucket.
func newBucket[K comparable, V any](hmap map[K]V) *bucket[K, V] {
	return &bucket[K, V]{
		hmap: hmap,
	}
}
