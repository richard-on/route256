// Package workerpool implements a goroutine worker pool.
package workerpool

import (
	"context"
	"sync"
)

// WorkerFunc represents a function used as a callback in the worker pool.
type WorkerFunc[In, Out any] func(ctx context.Context, arg In) (Out, error)

// Task represents an incoming worker pool task.
type Task[In, Out any] struct {
	// fn is the callback function which will be used in the worker pool.
	fn WorkerFunc[In, Out]
	// arg is the argument passed to fn.
	arg In
}

// Result represents the outcome of the callback function in the worker pool.
type Result[Out any] struct {
	// Value is the main result of the function.
	Value Out
	// Err is the error that could have happened in a callback.
	Err error
}

// Pool represents a worker pool with fixed number of concurrently running goroutines.
type Pool[In, Out any] struct {
	// maxWorkers is the maximum number of concurrent goroutines in the pool.
	maxWorkers int
	// taskGetChan is a channel that is used as a queue for incoming tasks.
	taskGetChan chan Task[In, Out]
	// taskProcessChan is a channel that is used as a queue for active worker.
	taskProcessChan chan Task[In, Out]
	// resChan is a channel for outgoing results from workers.
	resChan chan Result[Out]
	// results is the slice of Result of the worker pool.
	results []Result[Out]
	// closeSignal is a channel used to signal the Pool shutdown.
	closeSignal chan struct{}
	// accumulateSignal is a channel used to signal the end of result accumulation.
	accumulateSignal chan struct{}
	closeOnce        sync.Once
	// wg is a WaitGroup used to assure all workers have exited.
	wg sync.WaitGroup
}

// New creates a new Pool instance.
func New[In, Out any](ctx context.Context, maxWorkers int) *Pool[In, Out] {
	// maxWorkers should be at least 1.
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	p := &Pool[In, Out]{
		maxWorkers:       maxWorkers,
		taskGetChan:      make(chan Task[In, Out]),
		taskProcessChan:  make(chan Task[In, Out]),
		resChan:          make(chan Result[Out]),
		closeSignal:      make(chan struct{}),
		accumulateSignal: make(chan struct{}),
		wg:               sync.WaitGroup{},
	}

	// Start the pool.
	go p.run(ctx)

	return p
}

// SubmitOne adds one new Task to the pool.
// This operation may be blocking.
func (p *Pool[In, Out]) SubmitOne(fn WorkerFunc[In, Out], task In) {
	p.taskGetChan <- Task[In, Out]{fn, task}
}

// SubmitMany adds a slice of new Tasks with the same WorkerFunc to the pool.
// This operation may be blocking.
func (p *Pool[In, Out]) SubmitMany(fn WorkerFunc[In, Out], tasks []In) {
	for i := range tasks {
		p.taskGetChan <- Task[In, Out]{fn, tasks[i]}
	}
}

// GetResult retrieves a slice of Result.
// GetResult should generally be called after Wait or StopNow.
func (p *Pool[In, Out]) GetResult() []Result[Out] {
	return p.results
}

// Wait stops accepting new tasks to the Pool and
// waits for all tasks to execute, then stops the Pool.
func (p *Pool[In, Out]) Wait() {
	p.closeOnce.Do(func() {
		close(p.taskGetChan)
	})

	<-p.closeSignal
}

// StopNow stops accepting new tasks to the Pool
// while not waiting for already active or enqueued tasks.
func (p *Pool[In, Out]) StopNow() {
	p.closeOnce.Do(func() {
		close(p.taskGetChan)
	})
}

// run starts the pool and handles its shutdown.
func (p *Pool[In, Out]) run(ctx context.Context) {
	// close the channel to signal that the Pool have stopped running.
	defer close(p.closeSignal)

	workerCount := 0

	// Write results from resChan to results in a separate goroutine.
	go p.accumulateResult()

	// Spawn maxWorkers number of goroutines.
	for workerCount < p.maxWorkers {
		p.wg.Add(1)
		go worker(ctx, &p.wg, p.taskProcessChan, p.resChan)
		workerCount++
	}

	for {
		// Try to retrieve a new task.
		task, ok := <-p.taskGetChan
		if !ok {
			// If the channel is closed, no more tasks could be processed,
			// so initiate shutdown.
			close(p.taskProcessChan)
			break
		}

		// Send the task to processing channel.
		p.taskProcessChan <- task
	}

	// Wait for all workers to exit.
	p.wg.Wait()
	// Close the channel to signal to accumulator that no more results will arrive.
	close(p.resChan)
	// Wait for all results to be written.
	<-p.accumulateSignal
}

// accumulateResult writes Result sent to resChan to results slice.
func (p *Pool[In, Out]) accumulateResult() {
	// close the channel to signal that all results have been written.
	defer close(p.accumulateSignal)
	for res := range p.resChan {
		p.results = append(p.results, res)
	}
}

func worker[In, Out any](ctx context.Context, wg *sync.WaitGroup,
	taskChan <-chan Task[In, Out], resChan chan<- Result[Out]) {

	for task := range taskChan {
		select {
		// if context was cancelled, don't execute the task.
		case <-ctx.Done():
		default:
			// Otherwise, execute task.
			res, err := task.fn(ctx, task.arg)
			// And write its result to resChan.
			resChan <- Result[Out]{res, err}
		}
	}

	// Signal that this worker has exited.
	wg.Done()
	return
}
