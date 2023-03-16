// Package workerpool implements a goroutine worker pool.
//
// This package is heavily inspired by github.com/gammazero/workerpool.
// The main difference is the requirement to pass context as an argument when creating a pool.
// Therefore, this implementation can stop executing tasks as soon as context is cancelled.
package workerpool

import (
	"context"
	"sync"
	"time"

	"gitlab.ozon.dev/rragusskiy/homework-1/lib/queue"
)

const (
	// idleTimeout is the minimum period of time that worker needs to be idle to be stopped.
	idleTimeout = 2 * time.Second
)

// Pool represents a worker pool.
type Pool struct {
	// maxWorkers is the maximum number of concurrent goroutines in the pool.
	maxWorkers int
	// taskChan is a channel that is used as a queue for incoming tasks.
	taskChan chan func()
	// workerChan is a channel that is used as a queue for active worker.
	workerChan chan func()
	// waitQueue is the queue.Queue of enqueued tasks.
	waitQueue queue.Queue[func()]
	// stopChan is a channel used to signal the Pool shutdown.
	stopChan chan struct{}
	// stopOnce is used to ensure Pool is stopped only once.
	stopOnce sync.Once
	// waitEnqueued indicates whether pool needs to wait
	// for all tasks in waitQueue to execute before stopping.
	waitEnqueued bool
	// wg is a WaitGroup used to assure all workers have exited.
	wg sync.WaitGroup
}

// New creates a new Pool instance.
func New(ctx context.Context, maxWorkers int) *Pool {
	// maxWorkers should be at least 1.
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	pool := &Pool{
		maxWorkers: maxWorkers,
		taskChan:   make(chan func()),
		workerChan: make(chan func()),
		stopChan:   make(chan struct{}),
	}

	// Start the pool.
	go pool.run(ctx)

	return pool
}

// Submit adds task to a pool. This operation will not be blocking.
// Task might be executed immediately or enqueued.
func (p *Pool) Submit(task func()) {
	if task != nil {
		p.taskChan <- task
	}
}

// Stop stops accepting new tasks to the Pool
// while only waiting for currently running tasks to execute.
// It will not execute enqueued tasks.
func (p *Pool) Stop() {
	p.stop(false)
}

// Wait stops accepting new tasks to the Pool and
// waits for all tasks to execute, including enqueued ones and then stops the Pool.
func (p *Pool) Wait() {
	p.stop(true)
}

// stop tells the runner to shut down, and whether to complete enqueued tasks.
func (p *Pool) stop(wait bool) {
	// Task channel should only be closed once.
	p.stopOnce.Do(func() {
		p.waitEnqueued = wait
		// Close task channel so that no more tasks can be accepted.
		// Wait for currently running tasks to finish.
		close(p.taskChan)
	})

	// Send a signal to stop.
	<-p.stopChan
}

// run starts the pool and handles its shutdown.
func (p *Pool) run(ctx context.Context) {
	defer close(p.stopChan)

	timeout := time.NewTimer(idleTimeout)

	// Start processing tasks.
	workerCount := p.processTasks(ctx, timeout)

	if p.waitEnqueued {
		// Run enqueued tasks.
		for p.waitQueue.Len() != 0 {
			// A worker is ready, so assign a new task to it.
			p.workerChan <- p.waitQueue.Pop()
		}
	}

	// Send exit signal to all still active workers.
	for workerCount > 0 {
		p.workerChan <- nil
		workerCount--
	}

	// Wait for all workers to exit.
	p.wg.Wait()
	timeout.Stop()
}

// processTasks manages Pool tasks and workers.
// It assigns tasks to available worker or enqueues them.
func (p *Pool) processTasks(ctx context.Context, timeout *time.Timer) int {
	workerCount := 0
	idle := false
	for {
		// Check tasks in waitQueue.
		if p.waitQueue.Len() != 0 {
			if !p.processWaitQueue() {
				// If taskChan is closed, no more tasks could be processed.
				return workerCount
			}
			continue
		}

		select {
		case task, ok := <-p.taskChan:
			if !ok {
				// If taskChan is closed, no more tasks could be processed.
				return workerCount
			}

			select {
			case p.workerChan <- task: // If available worker exists, assign task to it.
			default:
				// If workerCount did not hit the limit,
				// spawn a new worker and assign current task to it.
				if workerCount < p.maxWorkers {
					p.wg.Add(1)
					go worker(ctx, task, p.workerChan, &p.wg)
					workerCount++
				} else {
					p.waitQueue.Push(task) // If the limit was hit, enqueue current task.
				}
			}
			idle = false

		case <-timeout.C:
			if idle && workerCount > 0 {
				if p.killIdleWorker() {
					workerCount--
				}
			}
			idle = true
			timeout.Reset(idleTimeout)
		}
	}
}

// processWaitQueue manages the wait queue.
// It enqueues tasks and dequeues them if a free worker is found.
func (p *Pool) processWaitQueue() bool {
	select {
	case task, ok := <-p.taskChan:
		if !ok {
			// If taskChan is closed, no more tasks could be enqueued.
			return false
		}
		// Push current task to the end of the queue.
		p.waitQueue.Push(task)

	// If a free worker is discovered, assign the first element on the queue to it.
	case p.workerChan <- p.waitQueue.Peek():
		// Then, pop this element from wait queue.
		p.waitQueue.Pop()
	}

	return true
}

// killIdleWorker tries to kill the worker.
func (p *Pool) killIdleWorker() bool {
	select {
	case p.workerChan <- nil:
		// Sent kill signal to worker.
		return true
	default:
		// No ready workers. All, if any, workers are busy.
		return false
	}
}

// worker runs the tasks.
// It will not execute provided task if context was cancelled.
func worker(ctx context.Context, task func(), workerChan chan func(), wg *sync.WaitGroup) {
	for task != nil {
		select {
		case <-ctx.Done():
			// if context was cancelled, don't execute the task.
		default:
			// Otherwise, execute task
			task()
		}

		// Get a new task.
		task = <-workerChan
	}
	wg.Done()
}
