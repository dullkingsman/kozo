package queue

import (
	"sync"
)

// Queue is a thread-safe FIFO data structure implemented with a circular buffer.
type Queue[T any] struct {
	mu    sync.Mutex
	data  []T
	head  int
	tail  int
	count int
}

// New returns a new empty Queue.
func New[T any]() *Queue[T] {
	return &Queue[T]{
		data: make([]T, 2), // Initial small capacity
	}
}

// NewWithCapacity returns a new empty Queue with pre-allocated capacity.
func NewWithCapacity[T any](capacity int) *Queue[T] {
	if capacity < 1 {
		capacity = 1
	}
	return &Queue[T]{
		data: make([]T, capacity),
	}
}

// Enqueue adds an element to the back of the queue.
func (q *Queue[T]) Enqueue(v T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.count == len(q.data) {
		q.resize()
	}

	q.data[q.tail] = v
	q.tail = (q.tail + 1) % len(q.data)
	q.count++
}

// Dequeue removes and returns the front element of the queue.
// Returns (zero-value, false) if the queue is empty.
func (q *Queue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.count == 0 {
		var zero T
		return zero, false
	}

	v := q.data[q.head]

	// Zero out the element to prevent memory leaks (GC can reclaim it)
	var zero T
	q.data[q.head] = zero

	q.head = (q.head + 1) % len(q.data)
	q.count--

	return v, true
}

// Peek returns the front element of the queue without removing it.
// Returns (zero-value, false) if the queue is empty.
func (q *Queue[T]) Peek() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.count == 0 {
		var zero T
		return zero, false
	}

	return q.data[q.head], true
}

// IsEmpty returns true if the queue has no elements.
func (q *Queue[T]) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.count == 0
}

// Len returns the current number of elements in the queue.
func (q *Queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.count
}

// Clear discards all elements from the queue.
func (q *Queue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Zero out all elements to assist GC
	var zero T
	for i := 0; i < len(q.data); i++ {
		q.data[i] = zero
	}

	q.head = 0
	q.tail = 0
	q.count = 0
}

// resize grows the underlying slice. Must be called with lock held.
func (q *Queue[T]) resize() {
	newCap := len(q.data) * 2
	if newCap == 0 {
		newCap = 1
	}
	newData := make([]T, newCap)

	for i := 0; i < q.count; i++ {
		newData[i] = q.data[(q.head+i)%len(q.data)]
	}

	q.data = newData
	q.head = 0
	q.tail = q.count
}
