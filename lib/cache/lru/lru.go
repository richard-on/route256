// Package lru implements LRU cache with TTL.
package lru

import (
	"sync"
	"time"

	"gitlab.ozon.dev/rragusskiy/homework-1/lib/hashmap"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/linkedlist"
)

// LRU is an LRU cache with TTL.
// It accepts comparable key and any value.
type LRU[K comparable, V any] struct {
	// hmap is a custom sharded hashmap.
	hmap *hashmap.HashMap[K, *linkedlist.Node[*Entry[K, V]]]
	// list is a doubly linked list.
	list *linkedlist.LinkedList[*Entry[K, V]]
	// capacity is a maximum number of elements in the cache.
	capacity int
	// ttl is a global time to live in this cache.
	ttl time.Duration

	mu sync.RWMutex
}

// Entry is a cache entry.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
	// expire is a time when this entry will be expired.
	expire time.Time
}

// New creates a new LRU cache.
func New[K comparable, V any](capacity int, ttl time.Duration) *LRU[K, V] {
	// We need at least one element in the cache.
	if capacity < 1 {
		capacity = 1
	}

	return &LRU[K, V]{
		hmap:     hashmap.New[K, *linkedlist.Node[*Entry[K, V]]](capacity),
		list:     linkedlist.New[*Entry[K, V]](),
		capacity: capacity,
		ttl:      ttl,
	}
}

// Cap returns a capacity of the cache.
func (l *LRU[K, V]) Cap() int {
	return l.capacity
}

// TTL returns a TTL of the cache.
func (l *LRU[K, V]) TTL() time.Duration {
	return l.ttl
}

// Get returns a value from the cache and
// a bool flag that signals whether key was found.
func (l *LRU[K, V]) Get(key K) (V, bool) {
	// Check and evict expired entries.
	l.evictExpired()

	// Lock the cache for reading.
	l.mu.RLock()
	defer l.mu.RUnlock()

	entry, ok := l.hmap.Get(key)
	if !ok {
		// Return zero value and false if key was not found.
		var res V
		return res, false
	}

	// Move the entry to the front of the list.
	l.list.MoveToFront(entry)
	return entry.Payload.Value, true
}

// Set sets a value in the cache.
func (l *LRU[K, V]) Set(key K, value V) {
	// Check and evict expired entries.
	l.evictExpired()

	l.mu.Lock()
	defer l.mu.Unlock()

	if entry, ok := l.hmap.Get(key); ok {
		// Update the value and move the entry to the front of the list.
		l.list.Remove(entry)
		newEntry := l.list.PushFront(&Entry[K, V]{
			Key:    key,
			Value:  value,
			expire: time.Now().Add(l.ttl),
		})
		// Also update the entry in the hashmap.
		l.hmap.Set(key, newEntry)
	} else {
		// Check if we need to evict an entry.
		if l.hmap.Len() >= l.capacity {
			// Evict the least recently used entry (last in the list).
			node := l.list.Back()
			val := l.list.Remove(node)
			l.hmap.Delete(val.Key)
		}

		// Add a new entry to the front of the list and to the hashmap.
		newElem := l.list.PushFront(&Entry[K, V]{
			Key:    key,
			Value:  value,
			expire: time.Now().Add(l.ttl),
		})
		l.hmap.Set(key, newElem)
	}
}

// evictExpired evicts expired entries from the cache.
func (l *LRU[K, V]) evictExpired() {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Iterate over the list and evict expired entries.
	entry := l.list.Front()
	for entry != nil {
		next := entry.Next()
		if entry.Payload.expire.Before(time.Now()) {
			l.list.Remove(entry)
			l.hmap.Delete(entry.Payload.Key)
		}
		entry = next
	}
}
