// Package queue implements a FIFO queue with circular buffer.
//
// This package is based on github.com/gammazero/deque.
package queue

// minCapacity is the smallest capacity that queue may have.
// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
const minCapacity = 16

// Queue represents an instance of FIFO queue.
// It contains items of the type specified by the type argument.
type Queue[T any] struct {
	buf    []T // buf stores elements of queue. It acts as a circular buffer.
	head   int // head is the position of queue head in the buf.
	tail   int // tail is the position of queue tail in the buf.
	count  int // count is the current number of elements in the buf.
	minCap int // minCap is the minimum size of buf. It will not shrink below this value.
}

// New creates a new Queue, optionally setting the current and minimum capacity.
// The Queue instance operates on items of the type specified by the type argument.
//
// Any size values supplied are rounded up to the nearest power of 2.
func New[T any](size ...int) *Queue[T] {
	var capacity, minimum int
	if len(size) >= 1 {
		capacity = size[0]
		if len(size) >= 2 {
			minimum = size[1]
		}
	}

	minCap := minCapacity
	// Round up to the nearest power of 2.
	for minCap < minimum {
		minCap <<= 1
	}

	var buf []T
	if capacity != 0 {
		bufSize := minCap
		// Round up to the nearest power of 2.
		for bufSize < capacity {
			bufSize <<= 1
		}
		buf = make([]T, bufSize)
	}

	return &Queue[T]{
		buf:    buf,
		minCap: minCap,
	}
}

// Len returns the number of elements currently stored in the queue. If q is nil, Len() returns zero.
func (q *Queue[T]) Len() int {
	if q == nil {
		return 0
	}
	return q.count
}

// Cap returns the current capacity of the queue. If q is nil, Cap() returns zero.
func (q *Queue[T]) Cap() int {
	if q == nil {
		return 0
	}
	return len(q.buf)
}

// Push pushes a new element to the end of the Queue.
func (q *Queue[T]) Push(elem T) {
	if q.count == len(q.buf) {
		q.grow()
	}

	q.buf[q.tail] = elem
	// Calculate new tail position.
	q.tail = q.next(q.tail)
	q.count++
}

func (q *Queue[T]) Pop() T {
	if q.count <= 0 {
		panic("queue: Pop() called on empty queue")
	}

	val := q.buf[q.head]

	var nilVal T
	// Set popped value in buffer to a nil value.
	q.buf[q.head] = nilVal
	// Calculate new head position.
	q.head = q.next(q.head)
	q.count--
	q.shrink()

	return val
}

// Peek returns the element at the front of the queue. This is the element
// that would be returned by PopFront(). This call panics if the queue is
// empty.
func (q *Queue[T]) Peek() T {
	if q.count <= 0 {
		panic("queue: Peek() called on empty queue")
	}
	return q.buf[q.head]
}

// next returns the next buffer position wrapping around buffer.
func (q *Queue[T]) next(i int) int {
	return (i + 1) & (len(q.buf) - 1) // bitwise modulus
}

// grow resizes up the buffer if it is full.
func (q *Queue[T]) grow() {
	if len(q.buf) == 0 {
		if q.minCap == 0 {
			q.minCap = minCapacity
		}
		q.buf = make([]T, q.minCap)
		return
	}

	q.resize()
}

// shrink resizes down if the buffer is 1/4 full.
func (q *Queue[T]) shrink() {
	if len(q.buf) > q.minCap && (q.count<<2) == len(q.buf) {
		q.resize()
	}
}

// resize resizes the queue to fit exactly twice its current contents. This is
// used to grow the queue when it is full, and also to shrink it when it is
// only a quarter full.
func (q *Queue[T]) resize() {
	// set new buffer size to twice its current size
	newBuf := make([]T, q.count<<1)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}
