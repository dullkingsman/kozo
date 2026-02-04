# Stack

A thread-safe, generic LIFO (Last-In, First-Out) stack implementation for Go.

## Features

- **Generic**: Works with any type `T` using Go 1.18+ generics.
- **Thread-Safe**: Safe for concurrent use across multiple goroutines using `sync.Mutex`.
- **Memory Optimized**: 
    - Supports pre-allocation to minimize "grow-and-copy" overhead.
    - Explicitly zeros elements on `Pop` and `Clear` to prevent memory leaks and assist the Garbage Collector.
- **Performance Focused**: Minimalistic API designed for compiler inlining and high-throughput.

## Installation

```bash
go get kozo/pkg/stack
```

## Quick Start

```go

import "github.com/dullkingsman/kozo/stack"

func main() {
    // Create a new stack for integers
    s := stack.New[int]()

    // Push elements
    s.Push(10)
    s.Push(20)

    // Peek at the top element without removing it
    if v, ok := s.Peek(); ok {
        fmt.Println("Top element:", v) // 20
    }

    // Pop elements
    v, ok := s.Pop() // 20, true
    v, ok = s.Pop()  // 10, true
    v, ok = s.Pop()  // 0, false (empty)
}
```

## API Reference

### Construction

- `New[T any]() *Stack[T]`: Creates an empty stack.
- `NewWithCapacity[T any](capacity int) *Stack[T]`: Creates an empty stack with pre-allocated capacity. Recommended when the maximum size is known to avoid re-allocations.

### Core Operations

- `Push(v T)`: Adds an element to the top of the stack.
- `Pop() (T, bool)`: Removes and returns the top element. Returns `(zero-value, false)` if the stack is empty.
- `Peek() (T, bool)`: Returns the top element without removing it. Returns `(zero-value, false)` if the stack is empty.

### State Metadata

- `Len() int`: Returns the current number of elements.
- `IsEmpty() bool`: Returns `true` if the stack contains no elements.

### Utility Operations

- `Swap() bool`: Swaps the top two elements. Returns `false` if the stack has fewer than two elements.
- `Clear()`: Discards all elements from the stack and zeros the underlying memory to assist GC.

## Optimizations

This implementation addresses deep runtime optimizations:

### 1. Memory Management
- **Zeroing on Pop**: When an element is popped, the underlying slice index is set to the zero value of `T`. This is critical when `T` contains pointers, as it allows the GC to reclaim memory immediately rather than waiting for the entire slice to be deallocated.
- **Slice Re-slicing**: Uses the `s = s[:len(s)-1]` idiom for O(1) removals, which is the most efficient way to manage slice-based stacks in Go.

### 2. Concurrency
- **Thread-Safety**: All operations are protected by a `sync.Mutex`, ensuring that the stack remains consistent even when accessed by hundreds of goroutines simultaneously.

### 3. Compiler & Runtime Efficiency
- **Generics**: Avoids "boxing" values into `interface{}`, eliminating reflection overhead and allowing for a more compact memory layout.
- **Inlining**: Methods like `Push`, `Pop`, and `Peek` are kept simple to encourage the Go compiler to inline them, removing function call overhead on hot paths.

## Usage Examples

### Concurrency Example

```go
s := stack.New[string]()
var wg sync.WaitGroup

for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(val int) {
        defer wg.Done()
        s.Push(fmt.Sprintf("item-%d", val))
    }(i)
}

wg.Wait()
fmt.Printf("Stack size: %d\n", s.Len()) // 100
```

### Performance Optimization with Capacity

```go
// Pre-allocate for 1000 items to avoid internal slice growth
s := stack.NewWithCapacity[int](1000)

for i := 0; i < 1000; i++ {
    s.Push(i)
}
```
