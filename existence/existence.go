package existence

// ExistenceClaim represents a filter condition:
// either "the value must be in this set" (Contains=true)
// or "the value must NOT be in this set" (Contains=false).
type ExistenceClaim[T any] struct {
	Values   []T  `json:"in"`
	Contains bool `json:"contains"`
}

// In creates an inclusive ExistenceClaim.
// The filter will match if the value is present in the provided values.
func In[T any](values ...T) ExistenceClaim[T] {
	return ExistenceClaim[T]{
		Values:   values,
		Contains: true,
	}
}

// NotIn creates an exclusive ExistenceClaim.
// The filter will match if the value is NOT present in the provided values.
func NotIn[T any](values ...T) ExistenceClaim[T] {
	return ExistenceClaim[T]{
		Values:   values,
		Contains: false,
	}
}

// Check determines if a value satisfies the existence claim using a custom equality function.
func (e ExistenceClaim[T]) Check(val T, equals func(T, T) bool) bool {
	found := false
	for _, v := range e.Values {
		if equals(v, val) {
			found = true
			break
		}
	}
	return found == e.Contains
}

// CheckComparable determines if a value satisfies the existence claim for comparable types.
func CheckComparable[T comparable](e ExistenceClaim[T], val T) bool {
	found := false
	for _, v := range e.Values {
		if v == val {
			found = true
			break
		}
	}
	return found == e.Contains
}

// IsEmpty returns true if the Values slice is empty.
func (e ExistenceClaim[T]) IsEmpty() bool {
	return len(e.Values) == 0
}

// Len returns the number of values in the claim.
func (e ExistenceClaim[T]) Len() int {
	return len(e.Values)
}

// Negate returns a new ExistenceClaim with the Contains flag flipped.
func (e ExistenceClaim[T]) Negate() ExistenceClaim[T] {
	return ExistenceClaim[T]{
		Values:   e.Values,
		Contains: !e.Contains,
	}
}

// Apply filters a slice based on the existence claim.
func (e ExistenceClaim[T]) Apply(slice []T, equals func(T, T) bool) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if e.Check(v, equals) {
			result = append(result, v)
		}
	}
	return result
}
