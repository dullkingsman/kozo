package _range

import "cmp"

// Range represents an interval.
type Range[T any] struct {
	Min *RangeItem[T] `json:"min"`
	Max *RangeItem[T] `json:"max"`
}

// RangeItem represents a boundary of an range.
type RangeItem[T any] struct {
	Value     *T   `json:"value"`
	Inclusive bool `json:"inclusive"`
}

// New creates a new Range with the given boundaries.
func New[T any](min, max *RangeItem[T]) Range[T] {
	return Range[T]{
		Min: min,
		Max: max,
	}
}

// Closed creates an inclusive range [min, max].
func Closed[T any](min, max T) Range[T] {
	return Range[T]{
		Min: &RangeItem[T]{Value: &min, Inclusive: true},
		Max: &RangeItem[T]{Value: &max, Inclusive: true},
	}
}

// Open creates an exclusive range (min, max).
func Open[T any](min, max T) Range[T] {
	return Range[T]{
		Min: &RangeItem[T]{Value: &min, Inclusive: false},
		Max: &RangeItem[T]{Value: &max, Inclusive: false},
	}
}

// HalfOpen creates a half-open range [min, max).
func HalfOpen[T any](min, max T) Range[T] {
	return Range[T]{
		Min: &RangeItem[T]{Value: &min, Inclusive: true},
		Max: &RangeItem[T]{Value: &max, Inclusive: false},
	}
}

// GreaterThan creates an exclusive range (min, +inf).
func GreaterThan[T any](min T) Range[T] {
	return Range[T]{
		Min: &RangeItem[T]{Value: &min, Inclusive: false},
	}
}

// AtLeast creates an inclusive range [min, +inf).
func AtLeast[T any](min T) Range[T] {
	return Range[T]{
		Min: &RangeItem[T]{Value: &min, Inclusive: true},
	}
}

// LessThan creates an exclusive range (-inf, max).
func LessThan[T any](max T) Range[T] {
	return Range[T]{
		Max: &RangeItem[T]{Value: &max, Inclusive: false},
	}
}

// AtMost creates an inclusive range (-inf, max].
func AtMost[T any](max T) Range[T] {
	return Range[T]{
		Max: &RangeItem[T]{Value: &max, Inclusive: true},
	}
}

// Contains determines if a value falls within the range using a custom less function.
func (r Range[T]) Contains(val T, less func(T, T) bool) bool {
	if r.Min != nil && r.Min.Value != nil {
		min := *r.Min.Value
		if r.Min.Inclusive {
			// val < min
			if less(val, min) {
				return false
			}
		} else {
			// val <= min  => !(min < val)
			if !less(min, val) {
				return false
			}
		}
	}

	if r.Max != nil && r.Max.Value != nil {
		max := *r.Max.Value
		if r.Max.Inclusive {
			// val > max
			if less(max, val) {
				return false
			}
		} else {
			// val >= max => !(val < max)
			if !less(val, max) {
				return false
			}
		}
	}

	return true
}

// ContainsOrdered determines if a value falls within the range for ordered types.
func ContainsOrdered[T cmp.Ordered](r Range[T], val T) bool {
	return r.Contains(val, func(a, b T) bool {
		return a < b
	})
}

// IsBounded returns true if both min and max are set.
func (r Range[T]) IsBounded() bool {
	return r.Min != nil && r.Min.Value != nil && r.Max != nil && r.Max.Value != nil
}

// IsAny returns true if neither min nor max are set (matches everything).
func (r Range[T]) IsAny() bool {
	minUnbounded := r.Min == nil || r.Min.Value == nil
	maxUnbounded := r.Max == nil || r.Max.Value == nil
	return minUnbounded && maxUnbounded
}
