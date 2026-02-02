package set

import (
	"sync"
)

// AnySet is a thread-safe set for any type T, using a custom equality function.
// Since it doesn't require T to be comparable, it uses a slice internally,
// making core operations O(n).
type AnySet[T any] struct {
	mu     sync.RWMutex
	items  []T
	equals func(T, T) bool
}

// NewAny creates a new AnySet for any type T, using the provided equality function.
func NewAny[T any](equals func(T, T) bool, items ...T) *AnySet[T] {
	s := &AnySet[T]{
		items:  make([]T, 0, len(items)),
		equals: equals,
	}
	s.Add(items...)
	return s
}

// Add adds one or more items to the set.
func (s *AnySet[T]) Add(items ...T) {
	if len(items) == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, item := range items {
		if !s.containsUnsafe(item) {
			s.items = append(s.items, item)
		}
	}
}

// Remove removes one or more items from the set.
func (s *AnySet[T]) Remove(items ...T) {
	if len(items) == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, item := range items {
		for i, existing := range s.items {
			if s.equals(existing, item) {
				// Efficiently remove by swapping with last element
				l := len(s.items)
				s.items[i] = s.items[l-1]
				// Zero out to assist GC
				var zero T
				s.items[l-1] = zero
				s.items = s.items[:l-1]
				break
			}
		}
	}
}

// Contains returns true if the set contains the item.
func (s *AnySet[T]) Contains(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.containsUnsafe(item)
}

func (s *AnySet[T]) containsUnsafe(item T) bool {
	for _, existing := range s.items {
		if s.equals(existing, item) {
			return true
		}
	}
	return false
}

// Pop removes and returns an arbitrary item from the set.
// Returns (zero-value, false) if the set is empty.
func (s *AnySet[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	l := len(s.items)
	if l == 0 {
		var zero T
		return zero, false
	}

	index := l - 1
	item := s.items[index]

	// Zero out to assist GC
	var zero T
	s.items[index] = zero
	s.items = s.items[:index]

	return item, true
}

// Len returns the number of items in the set.
func (s *AnySet[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

// IsEmpty returns true if the set contains no items.
func (s *AnySet[T]) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items) == 0
}

// Clear removes all items from the set.
func (s *AnySet[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Zero out to assist GC
	var zero T
	for i := range s.items {
		s.items[i] = zero
	}
	s.items = s.items[:0]
}

// ToSlice returns a slice containing all items in the set.
func (s *AnySet[T]) ToSlice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]T, len(s.items))
	copy(res, s.items)
	return res
}

// Iter iterates over the items in the set and calls the provided function for each item.
// If the function returns false, iteration stops.
func (s *AnySet[T]) Iter(fn func(T) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.items {
		if !fn(item) {
			break
		}
	}
}

// Clone returns a new AnySet with the same items.
func (s *AnySet[T]) Clone() *AnySet[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := &AnySet[T]{
		items:  make([]T, len(s.items)),
		equals: s.equals,
	}
	copy(res.items, s.items)
	return res
}

// Union returns a new set containing all items from both sets.
func (s *AnySet[T]) Union(other *AnySet[T]) *AnySet[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	res := &AnySet[T]{
		items:  make([]T, 0, len(s.items)+len(other.items)),
		equals: s.equals,
	}
	res.items = append(res.items, s.items...)
	for _, item := range other.items {
		if !res.containsUnsafe(item) {
			res.items = append(res.items, item)
		}
	}
	return res
}

// Intersect returns a new set containing only items present in both sets.
func (s *AnySet[T]) Intersect(other *AnySet[T]) *AnySet[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	res := &AnySet[T]{
		items:  make([]T, 0),
		equals: s.equals,
	}
	for _, item := range s.items {
		if other.containsUnsafe(item) {
			res.items = append(res.items, item)
		}
	}
	return res
}

// Difference returns a new set containing items present in s but not in other.
func (s *AnySet[T]) Difference(other *AnySet[T]) *AnySet[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	res := &AnySet[T]{
		items:  make([]T, 0),
		equals: s.equals,
	}
	for _, item := range s.items {
		if !other.containsUnsafe(item) {
			res.items = append(res.items, item)
		}
	}
	return res
}

// SymmetricDifference returns a new set containing items present in either s or other, but not both.
func (s *AnySet[T]) SymmetricDifference(other *AnySet[T]) *AnySet[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	res := &AnySet[T]{
		items:  make([]T, 0),
		equals: s.equals,
	}
	for _, item := range s.items {
		if !other.containsUnsafe(item) {
			res.items = append(res.items, item)
		}
	}
	for _, item := range other.items {
		if !s.containsUnsafe(item) {
			res.items = append(res.items, item)
		}
	}
	return res
}

// IsSubset returns true if all items in s are also in other.
func (s *AnySet[T]) IsSubset(other *AnySet[T]) bool {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	if len(s.items) > len(other.items) {
		return false
	}

	for _, item := range s.items {
		if !other.containsUnsafe(item) {
			return false
		}
	}
	return true
}

// IsSuperset returns true if all items in other are also in s.
func (s *AnySet[T]) IsSuperset(other *AnySet[T]) bool {
	return other.IsSubset(s)
}

// Equal returns true if both sets contain the same items.
func (s *AnySet[T]) Equal(other *AnySet[T]) bool {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	if len(s.items) != len(other.items) {
		return false
	}

	for _, item := range s.items {
		if !other.containsUnsafe(item) {
			return false
		}
	}
	return true
}
