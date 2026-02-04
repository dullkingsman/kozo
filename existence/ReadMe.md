# ExistenceClaim

`ExistenceClaim[T]` is a data structure used for filtering. It defines whether a value must be present in a specific set (`Contains=true`) or absent from it (`Contains=false`).

This is particularly useful for API filtering, where you might want to express logic like "status IN (active, pending)" or "id NOT IN (1, 2, 3)".

## Features

- **Generic**: Works with any type `T`.
- **JSON Ready**: Struct tags map `Values` to `in` for clean API payloads.
- **Flexible Equality**: Use `Check` with a custom equality function for complex types, or `CheckComparable` for primitive types.
- **Utility Methods**: Easily negate, check length, or apply the filter to entire slices.

## Installation

```bash
go get kozo/pkg/existence
```

## Basic Usage

### Construction

```go
import "github.com/dullkingsman/kozo/existence"

// Inclusive filter (In)
claim := existence.In("active", "pending")

// Exclusive filter (NotIn)
claim := existence.NotIn("deleted", "archived")
```

### Checking Values

```go
// For comparable types
match := existence.CheckComparable(claim, "active")

// For complex types
match := claim.Check(user, func(a, b User) bool {
    return a.ID == b.ID
})
```

### Applying to Slices

```go
users := []User{...}
activeUsers := claim.Apply(users, func(a, b User) bool {
    return a.Status == b
})
```

## API Reference

### Construction
- `In[T any](values ...T) ExistenceClaim[T]`: Creates a claim where values must be present.
- `NotIn[T any](values ...T) ExistenceClaim[T]`: Creates a claim where values must be absent.

### Operations
- `Check(val T, equals func(T, T) bool) bool`: Checks if a value satisfies the claim using a custom equality function.
- `CheckComparable[T comparable](e ExistenceClaim[T], val T) bool`: Optimized check for comparable types.
- `Apply(slice []T, equals func(T, T) bool) []T`: Returns a new slice containing only elements that satisfy the claim.
- `Negate() ExistenceClaim[T]`: Flips the `Contains` flag.
- `Len() int`: Returns the number of values in the claim.
- `IsEmpty() bool`: Returns true if the claim has no values.
