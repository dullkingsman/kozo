# Set

A collection of unique elements. This package provides two implementations of a thread-safe, generic set for Go.

## Implementations

### 1. `Set[T comparable]`
A high-performance set for types that support the `==` operator.
- **Performance**: $O(1)$ average for `Add`, `Remove`, and `Contains`.
- **Underlying Structure**: Uses Go's native `map[T]struct{}`.
- **Thread-Safety**: Protected by `sync.RWMutex` for concurrent reads.

### 2. `AnySet[T any]`
A flexible set for any type `T` that uses a custom equality function.
- **Performance**: $O(n)$ for `Add`, `Remove`, and `Contains`.
- **Underlying Structure**: Uses a slice `[]T`.
- **Use Case**: Best for small collections or complex types where native comparison is not possible.
- **Thread-Safety**: Protected by `sync.RWMutex`.

## Installation

```bash
go get kozo/pkg/set
```

## Quick Start

### Using `Set` (Comparable)

```go
import "github.com/dullkingsman/kozo/set"

func main() {
    // Create a new set of strings
    s := set.New("apple", "banana")

    s.Add("cherry")
    s.Contains("apple") // true
    s.Len()            // 3
    
    // Set operations
    other := set.New("banana", "date")
    union := s.Union(other) // apple, banana, cherry, date
}
```

### Using `AnySet` (With Equals Function)

```go
import "github.com/dullkingsman/kozo/set"

type User struct { ID int; Name string }

func main() {
    // Define equality by ID
    equals := func(a, b User) bool { return a.ID == b.ID }
    
    s := set.NewAny(equals)
    s.Add(User{ID: 1, Name: "Alice"})
    s.Add(User{ID: 1, Name: "Alice Redux"}) // Won't be added (duplicate ID)
    
    fmt.Println(s.Len()) // 1
}
```

## API Reference

The following methods are available on both `Set[T]` and `AnySet[T]`:

### Core Operations
- `Add(items ...T)`: Adds elements to the set.
- `Remove(items ...T)`: Removes elements from the set.
- `Contains(item T) bool`: Checks if an element exists.
- `Pop() (T, bool)`: Removes and returns an arbitrary element.
- `Clear()`: Discards all elements.

### State Metadata
- `Len() int`: Returns the number of elements.
- `IsEmpty() bool`: Returns `true` if empty.

### Set Operations (Returns new set)
- `Union(other)`: Elements in either set.
- `Intersect(other)`: Elements in both sets.
- `Difference(other)`: Elements in this set but not the other.
- `SymmetricDifference(other)`: Elements in either set but not both.

### Comparisons
- `IsSubset(other) bool`: This set is entirely contained in the other.
- `IsSuperset(other) bool`: This set contains all elements of the other.
- `Equal(other) bool`: Both sets contain exactly the same elements.

### Utility
- `ToSlice() []T`: Returns a slice of all elements.
- `Iter(func(T) bool)`: Iterates over elements. Return `false` to stop.
- `Clone()`: Returns a copy of the set.

## Optimizations

- **Memory Efficiency**: `Set[T]` uses `struct{}` as map values to minimize memory footprint.
- **Read Scalability**: Uses `sync.RWMutex` to allow multiple concurrent readers without blocking.
- **Batch Processing**: Variadic `Add` and `Remove` methods reduce lock contention for multiple items.
- **Zeroing**: `Pop` and `Remove` zero out deleted elements in `AnySet` to assist the Garbage Collector.
