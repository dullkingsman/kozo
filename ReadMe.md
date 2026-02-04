# kozo 

This is a selection of very few but very useful data structures 
I have used when working with go and are not provided by the 
standard library.

> The name is the phonetic reading of the Japanese word `構造` (Kōzō) meaning `structure`.

## Repository Structure

The library follows standard Go module conventions. And provides each data structure in its own package.

## Usage

### Stack

A generic, thread-safe stack implementation. See [Stack Documentation](stack/ReadMe.md) for detailed API and optimizations.

```go
import "github.com/dullkingsman/kozo/stack"

// Create a new stack
s := stack.New[int]()

// Push elements
s.Push(1)

// Pop elements
v, ok := s.Pop()
```

### Existence Claim

A generic `ExistenceClaim[T]` type for filtering, representing "In" or "NotIn" logic. See [Existence Claim Documentation](existence/ReadMe.md) for details.

```go
import "github.com/dullkingsman/kozo/existence"

// Values must be one of these
claim := existence.In("active", "pending")
match := existence.CheckComparable(claim, "active") // true
```

### Range

A generic `Range[T]` type for intervals. See [Range Documentation](range/ReadMe.md) for details.

```go
import "github.com/dullkingsman/kozo/_range"

// Inclusive range [10, 20]
r := _range.Closed(10, 20)
match := _range.ContainsOrdered(r, 15) // true
```

### Optional

A generic `Optional[T]` type that distinguishes between absent, null, and present values. See [Optional Documentation](optional/ReadMe.md) for detailed information on the three-state model and JSON support.

```go
import "github.com/dullkingsman/kozo/optional"

// Create a Some (present) optional with a value
o := optional.Some(42)

// Safe unwrapping
if v, ok := o.Unwrap(); ok {
    fmt.Println(v)
}
```

### Queue

A generic, thread-safe queue implementation optimized for $O(1)$ performance. See [Queue Documentation](queue/ReadMe.md) for detailed API and optimizations.

```go
import "github.com/dullkingsman/kozo/queue"

// Create a new queue
q := queue.New[int]()

// Enqueue elements
q.Enqueue(1)

// Dequeue elements
v, ok := q.Dequeue()
```

### Set

A generic, thread-safe set implementation. Provides `Set[T comparable]` for $O(1)$ performance and `AnySet[T any]` for custom equality. See [Set Documentation](set/ReadMe.md) for detailed API and optimizations.

```go
import "github.com/dullkingsman/kozo/set"

// Create a new set
s := set.New("apple", "banana")

// Add elements
s.Add("cherry")

// Check existence
exists := s.Contains("apple") // true
```
