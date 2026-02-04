# Optional

A generic `Optional[T]` type for Go that supports a three-state model, essential for distinguishing between "absent", "null", and "present with value".

## The Three-State Model

Unlike simple pointer-based or boolean-flagged optionals, this implementation distinguishes between three distinct states:

1.  **None**: The value is completely absent (e.g., a field missing from a JSON request).
2.  **Some(null)**: The value is present but explicitly set to null/nil (e.g., a JSON field set to `null`).
3.  **Some(value)**: The value is present and has a non-null value.

This is particularly useful for database updates where you need to know if you should:
-   Skip updating a column (None).
-   Set a column to NULL (Some null).
-   Update a column to a specific value (Some value).

## Installation

```bash
go get "kozo/pkg/optional"
```

## Basic Usage

### Construction

```go
import "github.com/dullkingsman/kozo/pkg/optional"

// Create a None (absent) optional
o1 := optional.None[int]()

// Create a Some (present) optional with a value
o2 := optional.Some(42)

// Zero-value is also None
var o3 optional.Optional[string]
```

### Checking State

```go
o := optional.Some(42)

o.IsSome()      // true
o.IsNone()      // false
o.IsNotNull()   // true
o.IsNull()      // false
o.IsNullOrNone() // false
```

### Accessing Values

```go
o := optional.Some(42)

// Safe unwrapping
if v, ok := o.Unwrap(); ok {
    fmt.Println(v)
}

// With default value
v := o.UnwrapOr(99)

// With lazy default
v := o.UnwrapOrElse(func() int { return computeDefault() })

// Panic if absent or null (use with caution)
v := o.Expect("value should be present")
```

## Functional API

```go
o := optional.Some(42)

// Filtering
even := o.Filter(func(v int) bool { return v%2 == 0 })

// Pattern Matching
o.Match(
    func(v int) { fmt.Printf("Got %d\n", v) },
    func()      { fmt.Println("Nothing") },
)

// Combining
combined := o.Or(optional.Some(99))
```

## JSON Support

The `Optional[T]` type implements `json.Marshaler` and `json.Unmarshaler`.

-   **None** marshals to `null` (but if used with `omitempty` in a struct, the field is omitted).
-   **Some(null)** marshals to `null`.
-   **Some(value)** marshals to the value itself.

Unmarshaling:
-   A missing field remains **None**.
-   A `null` value becomes **Some(null)**.
-   A concrete value becomes **Some(value)**.

```go
type User struct {
    Name optional.Optional[string] `json:"name,omitempty"`
}
```

## Advanced Operations

-   `Take(*Optional[T])`: Returns the value of an optional and leaves it as `None`.
-   `Clone()`: Creates a shallow copy of the optional.
-   `Xor(other)`: Returns an optional if exactly one of them is present.
