package data_structures

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Optional - Optional[T] represents an optional value of type T.
// It can either contain a value (Some) or be empty (None).
//
// It also further divides Some into Some(value) and Some(nil). This distinction is critical for database operations where
// "field not updated" vs "field set to null" have different meanings.
type Optional[T any] struct {
	value    *T
	nonEmpty bool
}

func (o Optional[T]) String() string {
	if o.IsSome() {
		if o.IsNull() {
			return "Some(null)"
		}

		return fmt.Sprintf("Some(%v)", *o.value)
	}

	return "None"
}

// =========================
// Zero Definition
// =========================

// IsZero returns true if the Optional is empty.
func (o Optional[T]) IsZero() bool {
	return !o.nonEmpty
}

// =========================
// Construction
// =========================

// Some creates an Optional containing the given value.
func Some[T any](v T) Optional[T] {
	return Optional[T]{value: &v, nonEmpty: true}
}

// None creates an empty Optional of type T.
func None[T any]() Optional[T] {
	return Optional[T]{value: nil, nonEmpty: false}
}

// TODO: Will be added based on need. It does not serve any particular purpose at the moment.
//// Zip combines two Options into one Optional of a tuple if both are Some.
//func Zip[T, U any](a Optional[T], b Optional[U]) Optional[Pair[*T, *U]] {
//	if a.IsSome() && b.IsSome() {
//		return Some(Pair[*T, *U]{First: a.value, Second: b.value})
//	}
//
//	return None[Pair[*T, *U]]()
//}

// =========================
// JSON Marshalling
// =========================

// MarshalJSON converts Optional[T] to JSON.
// - None → gets caught by standard JSON marshalling because of the omitzero tag since this will only be run after go 1.24
// - Some(value) → normal JSON of value
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}

// UnmarshalJSON converts JSON into Optional[T].
// - Missing field → None (handled by standard JSON unmarshalling)
// - JSON null → Some(nil)
// - JSON value → Some(value)
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)

	// Handle JSON null → Some(nil)
	if bytes.Equal(data, []byte("null")) {
		o.nonEmpty = true
		o.value = nil

		return nil
	}

	// Attempt to unmarshal normal value → Some(value)
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("cannot unmarshal Optional: %w", err)
	}

	o.value = &v
	o.nonEmpty = true

	return nil
}

// =========================
// Inspection
// =========================

// ContainsComparable returns true if Optional contains value (for comparable types).
func ContainsComparable[T comparable](o Optional[T], value T) bool {
	if o.IsNotNull() {
		return value == *o.value
	}

	return false
}

// Contains returns true if the Optional contains the specified value.
func (o Optional[T]) Contains(value T, compareFunc func(current T, incoming T) bool) bool {
	if o.IsNotNull() {
		return compareFunc(*o.value, value)
	}

	return false
}

// IsNullOrNone returns true if the Optional is empty or its value is null.
func (o Optional[T]) IsNullOrNone() bool {
	return o.IsNull() || o.IsNone()
}

// IsNotNull returns true if the Optional's value is not null.
func (o Optional[T]) IsNotNull() bool {
	return o.IsSome() && o.value != nil
}

// IsNull returns true if the Optional is not empty and its value is null.
func (o Optional[T]) IsNull() bool {
	return o.IsSome() && o.value == nil
}

// IsSome returns true if the Optional contains a value.
func (o Optional[T]) IsSome() bool {
	return !o.IsNone()
}

// IsNone returns true if the Optional is empty.
func (o Optional[T]) IsNone() bool {
	return o.nonEmpty == false
}

// =========================
// Access
// =========================

// ExpectPtr panics if the Optional is empty.
func (o Optional[T]) ExpectPtr(message string) *T {
	if o.IsSome() {
		return o.value
	}

	panic(message)
}

// Expect panics if the Optional is empty or its value is null.
func (o Optional[T]) Expect(message string) T {
	if o.IsNotNull() {
		return *o.value
	}

	panic(message)
}

// UnwrapPtr returns the value and true if present, otherwise nil and false.
func (o Optional[T]) UnwrapPtr() (*T, bool) {
	if o.IsSome() {
		return o.value, true
	}

	return nil, false
}

// Unwrap returns the value and true if present, otherwise zero value and false.
func (o Optional[T]) Unwrap() (T, bool) {
	if o.IsNotNull() {
		return *o.value, true
	}

	var zero T

	return zero, false
}

// UnwrapOrPtr returns the value or a default if empty or null.
//
// Note: This function returns nil when Some(nil), not the default value.
func (o Optional[T]) UnwrapOrPtr(defaultValue *T) *T {
	if o.IsSome() {
		return o.value
	}

	return defaultValue
}

// UnwrapOr returns the value or a default if empty or null.
func (o Optional[T]) UnwrapOr(defaultValue T) T {
	if o.IsNotNull() {
		return *o.value
	}

	return defaultValue
}

// UnwrapOrElsePtr returns the value or computes a default if empty.
func (o Optional[T]) UnwrapOrElsePtr(defaultFunc func() *T) *T {
	if o.IsSome() {
		return o.value
	}

	return defaultFunc()
}

// UnwrapOrElse returns the value or computes a default if empty.
func (o Optional[T]) UnwrapOrElse(defaultFunc func() T) T {
	if o.IsNotNull() {
		return *o.value
	}

	return defaultFunc()
}

// Take consumes the Optional and returns its value, leaving None behind.
func Take[T any](o *Optional[T]) Optional[T] {
	if o == nil {
		return None[T]()
	}

	var old = *o
	o.value = nil
	o.nonEmpty = false
	return old
}

// =========================
// Transformation
// =========================

// TODO: SINCE METHODS CAN NOT HAVE TYPE PARAMETER DEFINITIONS, TRANSFORMATION METHODS ARE NOT POSSIBLE TO IMPLEMENT THIS WAY.
// TODO: HOWEVER, WE CAN IMPLEMENT THEM AS EXTERNAL FUNCTIONS.
//// MapPtr applies a function to the value if the optional is not empty, returning a new Optional.
//func (o Optional[T]) MapPtr(f func(*T) *T) Optional[T] {
//	if o.IsSome() {
//		var zero T
//
//		var (
//			s = Some(zero)
//			r = f(o.value)
//		)
//
//		s.value = r
//		s.nonEmpty = true
//
//		return s
//	}
//
//	return None[T]()
//}
//
//// Map applies a function to the value if the optional is not empty and its value is not null, returning a new Optional.
//func (o Optional[T]) Map(f func(T) T) Optional[T] {
//	if o.IsNotNull() {
//		var r = f(*o.value)
//		return Some(r)
//	}
//
//	return None[T]()
//}
//
//// MapOrPtr applies a function to the value if present, else returns default.
//func (o Optional[T]) MapOrPtr(defaultValue T, f func(*T) *T) *T {
//	if o.IsSome() {
//		return f(o.value)
//	}
//
//	return &defaultValue
//}
//
//// MapOr applies a function to the value if present and its value is not null, else returns default.
//func (o Optional[T]) MapOr(defaultValue T, f func(T) T) T {
//	if o.IsNotNull() {
//		return f(*o.value)
//	}
//
//	return defaultValue
//}
//
//// MapOrElsePtr applies a function to the value if present, else computes a default.
//func (o Optional[T]) MapOrElsePtr(defaultFunc func() *T, f func(*T) *T) *T {
//	if o.IsSome() {
//		return f(o.value)
//	}
//
//	return defaultFunc()
//}
//
//// MapOrElse applies a function to the value if present and its value is not null, else computes a default.
//func (o Optional[T]) MapOrElse(defaultFunc func() T, f func(T) T) T {
//	if o.IsNotNull() {
//		return f(*o.value)
//	}
//
//	return defaultFunc()
//}
//
//// AndThenPtr chains another Optional-returning function if value is present, otherwise returns None.
//func (o Optional[T]) AndThenPtr(f func(*T) Optional[T]) Optional[T] {
//	if o.IsSome() {
//		return f(o.value)
//	}
//
//	return None[T]()
//}
//
//// AndThen chains another Optional-returning function if value is present and its value is not null, otherwise returns None.
//func (o Optional[T]) AndThen(f func(T) Optional[T]) Optional[T] {
//	if o.IsNotNull() {
//		return f(*o.value)
//	}
//
//	return None[T]()
//}

// =========================
// Copy
// =========================

// Clone creates a deep copy of the Optional.
//
// Note: For pointer or reference types (slices, maps), only the reference of the underlying value is copied.
func (o Optional[T]) Clone() Optional[T] {
	if o.IsNone() {
		return o
	}

	var n = Optional[T]{nonEmpty: o.nonEmpty}

	if o.IsNotNull() {
		var ptr = *o.value
		n.value = &ptr
	}

	return n
}

// =========================
// Filtering
// =========================

// FilterPtr returns Some the optional is not empty and the value satisfies the predicate, else None.
// predicate may receive a null value if the optional value is null.
func (o Optional[T]) FilterPtr(predicate func(*T) bool) Optional[T] {
	if o.IsSome() && predicate(o.value) {
		return o
	}

	return None[T]()
}

// Filter returns Some the optional is not empty, the value is not null and the value satisfies the predicate, else None.
func (o Optional[T]) Filter(predicate func(T) bool) Optional[T] {
	if o.IsNotNull() && predicate(*o.value) {
		return o
	}

	return None[T]()
}

// =========================
// Matching
// =========================

// MatchPtr applies the corresponding function to the contained value or a default function if the optional is empty.
func (o Optional[T]) MatchPtr(someFunc func(*T), noneFunc func()) {
	if o.IsSome() {
		someFunc(o.value)
	} else {
		noneFunc()
	}
}

// Match applies the corresponding function to the contained value if the value is not null or a default function if the optional is empty or the value is null.
func (o Optional[T]) Match(someFunc func(T), noneFunc func()) {
	if o.IsNotNull() {
		someFunc(*o.value)
	} else {
		noneFunc()
	}
}

// =========================
// Combining Options
// =========================

// Or returns self if Some, else other.
func (o Optional[T]) Or(other Optional[T]) Optional[T] {
	if o.IsSome() {
		return o
	}

	return other
}

// OrElse returns self if Some, else computes alternative Optional.
func (o Optional[T]) OrElse(f func() Optional[T]) Optional[T] {
	if o.IsSome() {
		return o
	}

	return f()
}

// Xor returns Some if exactly one of selves or other is Some, else None.
func (o Optional[T]) Xor(other Optional[T]) Optional[T] {
	if o.IsSome() && other.IsNone() {
		return o
	}

	if o.IsNone() && other.IsSome() {
		return other
	}

	return None[T]()
}
