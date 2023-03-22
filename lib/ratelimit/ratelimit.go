// Package ratelimit implements a simple rate limiter that uses leaky bucket algorithm,
// which can dynamically resize general and burst limit based on the current number of replicas.
//
// This package is based on:
// https://github.com/uber-go/ratelimit/blob/main/ratelimit.go
// which, in turn, was inspired by:
// https://github.com/prashantv/go-bench/blob/master/ratelimit
//
// This package's main goal is to keep the API simple
// while allowing it to be seamlessly used in distributed workloads.
// I plan to add token bucket algorithm and, possibly,
// limiter with cache storage, like Redis.
//
// This package is work in progress, I expect radical changes in the API.
package ratelimit

import (
	"context"
	"time"

	"github.com/andres-erbsen/clock"
)

const (
	defaultTimeWindow     = time.Second
	defaultBurst      int = 5
)

// Limiter is the interface that is used to put rate limit on given process.
type Limiter interface {
	// Take applies the rate limiter when it is called.
	// It may block to match desired rate limit.
	Take() time.Time
}

// Clock is the interface that allows usage of custom clocks
// and helps with testing.
type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}

// Option is the interface used to configure the limiter.
type Option interface {
	// apply applies given option to a config.
	apply(*config)
}

// config stores configuration for the limiter.
type config struct {
	clock Clock
	// burst is the maximum number of accumulated resources that can be used for bursts.
	burst int
	// per is the limiter window size.
	per time.Duration
	// replica holds data used to resize the Limiter.
	replica replica
}

// New creates a new limiter with given rate and options.
//
// Note: context is required for Limiter with replicas. Otherwise,
// you may pass nil or context.Background().
func New(ctx context.Context, rate int, opts ...Option) Limiter {
	c := config{
		clock: clock.New(),
		burst: defaultBurst,
		per:   defaultTimeWindow,
	}

	// Apply options.
	for _, opt := range opts {
		opt.apply(&c)
	}

	// Create a new limiter.
	// I plan on adding other algorithms
	return newLeakyBucketLimiter(ctx, rate, c)
}

type clockOption struct {
	clock Clock
}

func (o clockOption) apply(c *config) {
	c.clock = o.clock
}

// WithClock configures the limiter to use custom clock implementation.
func WithClock(clock Clock) Option {
	return clockOption{
		clock: clock,
	}
}

type burstOption int

func (o burstOption) apply(c *config) {
	c.burst = int(o)
}

// WithoutBurst configures the limiter to be strict.
// With this option, it will not accumulate resources for further bursts.
var WithoutBurst Option = burstOption(0)

// WithBurst configures custom maximum burst.
func WithBurst(burst int) Option {
	return burstOption(burst)
}

type perOption time.Duration

func (p perOption) apply(c *config) {
	c.per = time.Duration(p)
}

// Per sets custom time window for a Limiter.
// The default window is 1 second.
func Per(per time.Duration) Option {
	return perOption(per)
}

// Updater is the interface used to update the number of replicas.
//
// Updater can be used, for example, to fetch the number of running
// pod replicas in a cluster.
type Updater interface {
	// GetReplicaCount returns the number of currently running
	// replicas of any entity, that uses this rate limiter.
	GetReplicaCount(ctx context.Context) int
}

// replica holds data used to dynamically resize rate limiter.
type replica struct {
	// updater updates replicaCount.
	updater Updater
	// updateRate is the rate at which Updater should get new replicaCount.
	//
	// If updateRate is 1 * time.Second, Updater will be called once every second.
	updateRate time.Duration
	// replicaCount is the number of currently active replicas of rate limited entity.
	replicaCount int
}

type replicaOption replica

func (o replicaOption) apply(c *config) {
	c.replica = replica(o)
}

// WithReplicas configures the limiter to dynamically resize rate limit
// based on the number of replicas of given entity.
//
// For example, if initial rate limit was set to 10 RPS, and then
// one more process have been started in a medium tracked by Updater,
// rate limit on both processes will be equal to 5.
func WithReplicas(updater Updater, updateRate time.Duration) Option {
	return replicaOption{
		updater:    updater,
		updateRate: updateRate,
	}
}
