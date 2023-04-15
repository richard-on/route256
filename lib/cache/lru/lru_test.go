package lru

import (
	"sync"
	"testing"
	"time"
)

func BenchmarkDefaultMap(b *testing.B) {
	m := make(map[int]int, b.N)

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m[i] = i
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value, found := m[i]
			if found {
				_ = value
			}
		}
	})
}

func BenchmarkSyncMap(b *testing.B) {
	var m sync.Map

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.Store(i, i)
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value, found := m.Load(i)
			if found {
				_ = value
			}
		}
	})
}

func BenchmarkLRU(b *testing.B) {
	c := New[int, int](b.N, time.Hour)

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Set(i, i)
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value, found := c.Get(i)
			if found {
				_ = value
			}
		}
	})
}
