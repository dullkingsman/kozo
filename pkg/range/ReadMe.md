# Range

A generic `Range[T]` data structure for representing interval with inclusive or exclusive boundaries.

## Features

- **Generic**: Works with any type `T` using Go 1.18+ generics.
- **Inclusive/Exclusive**: Supports both `[min, max]` (closed), `(min, max)` (open), and `[min, max)` (half-open) ranges.
- **Unbounded**: Supports ranges like `> 10` (no upper bound) or `<= 20` (no lower bound).
- **JSON Ready**: Built-in JSON tags for seamless API integration.
- **Functional API**: Includes helper constructors and comparison methods.

## Installation

```bash
go get kozo/pkg/_range
```

## Quick Start

```go
import "kozo/pkg/_range"

func main() {
    // Create an inclusive range [10, 20]
    r := _range.Closed(10, 20)

    // Check if a value is within the range
    fmt.Println(_range.ContainsOrdered(r, 15)) // true
    fmt.Println(_range.ContainsOrdered(r, 25)) // false

    // Create a half-open range [10, 20)
    r2 := _range.HalfOpen(10, 20)
    fmt.Println(_range.ContainsOrdered(r2, 20)) // false
}
```

## API Reference

### Construction

- `Closed(min, max T)`: Creates `[min, max]`.
- `Open(min, max T)`: Creates `(min, max)`.
- `HalfOpen(min, max T)`: Creates `[min, max)`.
- `GreaterThan(min T)`: Creates `(min, +inf)`.
- `AtLeast(min T)`: Creates `[min, +inf)`.
- `LessThan(max T)`: Creates `(-inf, max)`.
- `AtMost(max T)`: Creates `(-inf, max]`.

### Verification

- `Contains(val T, less func(T, T) bool) bool`: Checks if `val` is in range using a custom comparison.
- `ContainsOrdered(r Range[T], val T) bool`: Optimized check for `cmp.Ordered` types.

### Metadata

- `IsBounded() bool`: Returns `true` if both `min` and `max` are set.
- `IsAny() bool`: Returns `true` if neither `min` nor `max` are set (matches everything).

## JSON Integration

The `Range[T]` struct uses pointer-based boundaries to represent unbounded states in JSON:

```json
{
  "min": { "value": 10, "inclusive": true },
  "max": null
}
```
The above represents `[10, +inf)`.
