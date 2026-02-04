# Queue

A thread-safe, generic FIFO (First-In, First-Out) queue implementation for Go, optimized for $O(1)$ performance.

## Features

- **Generic**: Works with any type `T` using Go 1.18+ generics.
- **Thread-Safe**: Safe for concurrent use across multiple goroutines using `sync.Mutex`.
- **$O(1)$ Performance**: Implemented with a circular buffer (ring buffer) to ensure that both `Enqueue` and `Dequeue` operations are $O(1)$ amortized, avoiding the $O(n)$ cost of shifting elements.
- **Memory Optimized**: 
    - Supports pre-allocation via `NewWithCapacity`.
    - Explicitly zeros elements on `Dequeue` and `Clear` to prevent memory leaks and assist the Garbage Collector.
- **Performance Focused**: Minimalistic API designed for compiler inlining and high-throughput.

## Installation

```bash
go get kozo/pkg/queue
```

## Quick Start

```go
import "github.com/dullkingsman/kozo/pkg/queue"

func main() {
    // Create a new queue for integers
    q := queue.New[int]()

    // Enqueue elements
    q.Enqueue(10)
    q.Enqueue(20)

    // Peek at the front element without removing it
    if v, ok := q.Peek(); ok {
        fmt.Println("Front element:", v) // 10
    }

    // Dequeue elements
    v, ok := q.Dequeue() // 10, true
    v, ok = q.Dequeue()  // 20, true
    v, ok = q.Dequeue()  // 0, false (empty)
}
```

## API Reference

### Construction

- `New[T any]() *Queue[T]`: Creates an empty queue.
- `NewWithCapacity[T any](capacity int) *Queue[T]`: Creates an empty queue with pre-allocated capacity. Recommended when the maximum size is known to avoid re-allocations.

### Core Operations

- `Enqueue(v T)`: Adds an element to the back of the queue.
- `Dequeue() (T, bool)`: Removes and returns the front element. Returns `(zero-value, false)` if the queue is empty.
- `Peek() (T, bool)`: Returns the front element without removing it. Returns `(zero-value, false)` if the queue is empty.

### State Metadata

- `Len() int`: Returns the current number of elements.
- `IsEmpty() bool`: Returns `true` if the queue contains no elements.

### Utility Operations

- `Clear()`: Discards all elements from the queue and zeros the underlying memory to assist GC.

## Optimizations

### 1. Circular Buffer Implementation
Most slice-based queues in Go use `s = s[1:]` for Dequeue, which is $O(1)$ but leaves the front of the underlying array unused. Eventually, this leads to memory bloat or $O(n)$ re-allocations. This implementation uses a circular buffer with `head` and `tail` pointers, ensuring true $O(1)$ operations and efficient memory reuse within the slice.

### 2. Memory Management
- **Zeroing on Dequeue**: When an element is dequeued, its slot in the underlying slice is set to the zero value of `T`. This is critical when `T` contains pointers, as it allows the GC to reclaim memory immediately.
- **Amortized Growth**: The queue grows exponentially when full, minimizing the number of allocations.

### 3. Concurrency
- **Thread-Safety**: All operations are protected by a `sync.Mutex`, making it safe for producer-consumer patterns across multiple goroutines.
