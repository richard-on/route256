package ratelimit

import (
	"context"
	"sync"
	"time"
)

// leakyLimiter is a struct that helps implement a leaky bucket rate limiting algorithm.
type leakyLimiter struct {
	clock Clock
	// last is the last time limiter has been taken.
	last time.Time
	// sleepFor is the current blocking duration.
	sleepFor time.Duration
	// per is the average rate limit leakyLimiter should meet.
	per time.Duration
	// maxBurst is the maximum allowed burst.
	maxBurst time.Duration
	// replica stores data regarding resizing.
	replica replica
	// mu is a mutex for sync.
	mu sync.Mutex
}

func newLeakyBucketLimiter(ctx context.Context, rate int, c config) Limiter {
	baseRate := c.per / time.Duration(rate)
	baseBurst := -1 * time.Duration(c.burst) * baseRate
	l := &leakyLimiter{
		clock:    c.clock,
		per:      baseRate,
		maxBurst: baseBurst,
		replica:  c.replica,
	}

	if l.replica.updater != nil {
		// Monitor the number of replicas and update rate limit in a separate goroutine.
		go l.monitorReplicaCount(ctx, baseRate, baseBurst)
	}

	return l
}

// Take applies the rate limiter when it is called.
// It may block to match desired rate limit.
func (l *leakyLimiter) Take() time.Time {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Get current timestamp.
	now := l.clock.Now()

	// If the limiter has never been taken before, allow the operation.
	if l.last.IsZero() {
		l.last = now
		return l.last
	}

	// Calculate the duration rate limiter should block.
	// This number may be negative as the last request could
	// have taken longer than the average rate limit.
	l.sleepFor += l.per - now.Sub(l.last)

	// Do not allow sleepFor to be below maxBurst as this
	// would cause a spike in rate after tasks speed up.
	if l.sleepFor < l.maxBurst {
		l.sleepFor = l.maxBurst
	}

	// If sleepFor is positive, block now.
	if l.sleepFor > 0 {
		l.clock.Sleep(l.sleepFor)
		l.sleepFor = 0
		l.last = now.Add(l.sleepFor)
	} else {
		l.last = now
	}

	return l.last
}

// monitorReplicaCount fetches the number of active limiters and adjust rate limit.
//
// If context is cancelled, synchronization is stopped immediately.
// The limiter however may continue working with the last applied limit.
func (l *leakyLimiter) monitorReplicaCount(ctx context.Context, baseRate, baseBurst time.Duration) {
	updateTicker := time.NewTicker(l.replica.updateRate)
	for {
		select {
		case <-ctx.Done():
			return
		// Once in an update rate, check replica count.
		case <-updateTicker.C:
			l.mu.Lock()
			// Get a new number of replicas.
			l.replica.replicaCount = l.replica.updater.GetReplicaCount(ctx)
			// Update rate limit. More replicas = slower individual rate.
			l.per = baseRate * time.Duration(l.replica.replicaCount)
			// Update maximum burst. More replicas = individual allowed burst should be smaller.
			l.maxBurst = baseBurst / time.Duration(l.replica.replicaCount)
			l.mu.Unlock()
		}
	}
}
