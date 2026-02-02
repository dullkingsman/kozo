package set

import (
	"sync"
)

// Set is a thread-safe, generic set for comparable types.
// It uses a map internally for O(1) average time complexity for core operations.
type Set[T comparable] struct {
	mu sync.RWMutex
	m  map[T]struct{}
}

// New creates a new Set for comparable types.
// If items are provided, they are added to the set.
func New[T comparable](items ...T) *Set[T] {
	s := &Set[T]{
		m: make(map[T]struct{}, len(items)),
	}
	s.Add(items...)
	return s
}

// Add adds one or more items to the set.
func (s *Set[T]) Add(items ...T) {
	if len(items) == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range items {
		s.m[item] = struct{}{}
	}
}

// Remove removes one or more items from the set.
func (s *Set[T]) Remove(items ...T) {
	if len(items) == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range items {
		delete(s.m, item)
	}
}

// Contains returns true if the set contains the item.
func (s *Set[T]) Contains(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.m[item]
	return ok
}

// Pop removes and returns an arbitrary item from the set.
// Returns (zero-value, false) if the set is empty.
func (s *Set[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for item := range s.m {
		delete(s.m, item)
		return item, true
	}

	var zero T
	return zero, false
}

// Len returns the number of items in the set.
func (s *Set[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.m)
}

// IsEmpty returns true if the set contains no items.
func (s *Set[T]) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.m) == 0
}

// Clear removes all items from the set.
func (s *Set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m = make(map[T]struct{})
}

// ToSlice returns a slice containing all items in the set.
// The order of items is non-deterministic.
func (s *Set[T]) ToSlice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]T, 0, len(s.m))
	for item := range s.m {
		res = append(res, item)
	}
	return res
}

// Iter iterates over the items in the set and calls the provided function for each item.
// If the function returns false, iteration stops.
func (s *Set[T]) Iter(fn func(T) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for item := range s.m {
		if !fn(item) {
			break
		}
	}
}

// Clone returns a new Set with the same items.
func (s *Set[T]) Clone() *Set[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := &Set[T]{
		m: make(map[T]struct{}, len(s.m)),
	}
	for item := range s.m {
		res.m[item] = struct{}{}
	}
	return res
}

// Union returns a new set containing all items from both sets.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	res := &Set[T]{
		m: make(map[T]struct{}, len(s.m)+len(other.m)),
	}
	for item := range s.m {
		res.m[item] = struct{}{}
	}
	for item := range other.m {
		res.m[item] = struct{}{}
	}
	return res
}

// Intersect returns a new set containing only items present in both sets.
func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	// Iterate over the smaller set for efficiency
	small, large := s, other
	if len(small.m) > len(large.m) {
		small, large = other, s
	}

	res := &Set[T]{
		m: make(map[T]struct{}),
	}
	for item := range small.m {
		if _, ok := large.m[item]; ok {
			res.m[item] = struct{}{}
		}
	}
	return res
}

// Difference returns a new set containing items present in s but not in other.
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	res := &Set[T]{
		m: make(map[T]struct{}),
	}
	for item := range s.m {
		if _, ok := other.m[item]; !ok {
			res.m[item] = struct{}{}
		}
	}
	return res
}

// SymmetricDifference returns a new set containing items present in either s or other, but not both.
func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	res := &Set[T]{
		m: make(map[T]struct{}),
	}
	for item := range s.m {
		if _, ok := other.m[item]; !ok {
			res.m[item] = struct{}{}
		}
	}
	for item := range other.m {
		if _, ok := s.m[item]; !ok {
			res.m[item] = struct{}{}
		}
	}
	return res
}

// IsSubset returns true if all items in s are also in other.
func (s *Set[T]) IsSubset(other *Set[T]) bool {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	if len(s.m) > len(other.m) {
		return false
	}

	for item := range s.m {
		if _, ok := other.m[item]; !ok {
			return false
		}
	}
	return true
}

// IsSuperset returns true if all items in other are also in s.
func (s *Set[T]) IsSuperset(other *Set[T]) bool {
	return other.IsSubset(s)
}

// Equal returns true if both sets contain the same items.
func (s *Set[T]) Equal(other *Set[T]) bool {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	if len(s.m) != len(other.m) {
		return false
	}

	for item := range s.m {
		if _, ok := other.m[item]; !ok {
			return false
		}
	}
	return true
}
