// Package hashmap implements a thread-safe hash map.
package hashmap

import (
	"sync/atomic"

	"github.com/dolthub/maphash"
)

// HashMap is a thread-safe hash map.
type HashMap[K comparable, V any] struct {
	// capacity is the max number of elements that can be stored in the map.
	capacity int
	// len is the number of elements in the map.
	// It is an atomic int64 to only use mutexes inside buckets.
	len atomic.Int64
	// hash is the hasher used to hash keys.
	// Maphash uses runtime hasher.
	hash maphash.Hasher[K]
	// buckets is the array of bucket.
	buckets []*bucket[K, V]
}

// New returns a new HashMap.
func New[K comparable, V any](capacity int) *HashMap[K, V] {
	if capacity < 1 {
		capacity = 1
	}

	s := HashMap[K, V]{
		capacity: capacity,
		hash:     maphash.NewHasher[K](),
		buckets:  make([]*bucket[K, V], capacity),
	}

	// Initialize buckets.
	for i := range s.buckets {
		s.buckets[i] = newBucket(make(map[K]V))
	}

	return &s
}

// Len returns the number of elements in the map.
func (m *HashMap[K, V]) Len() int {
	return int(m.len.Load())
}

// Get returns the value of the key from HashMap.
func (m *HashMap[K, V]) Get(key K) (V, bool) {
	// Get the bucket number.
	bucketNum := m.hash.Hash(key) & uint64(m.capacity-1)

	m.buckets[bucketNum].mu.RLock()
	defer m.buckets[bucketNum].mu.RUnlock()

	res, ok := m.buckets[bucketNum].hmap[key]
	return res, ok
}

// Set sets the value for the key-value pair in HashMap.
func (m *HashMap[K, V]) Set(key K, value V) {
	bucketNum := m.hash.Hash(key) & uint64(m.capacity-1)

	m.buckets[bucketNum].mu.Lock()
	defer m.buckets[bucketNum].mu.Unlock()

	if _, ok := m.buckets[bucketNum].hmap[key]; !ok {
		m.buckets[bucketNum].hmap[key] = value
		m.len.Add(1)
	}
}

// Delete deletes the key-value pair from HashMap.
func (m *HashMap[K, V]) Delete(key K) {
	bucketNum := m.hash.Hash(key) & uint64(m.capacity-1)

	m.buckets[bucketNum].mu.Lock()
	defer m.buckets[bucketNum].mu.Unlock()

	if _, ok := m.buckets[bucketNum].hmap[key]; ok {
		delete(m.buckets[bucketNum].hmap, key)
		m.len.Add(-1)
	}
}
