package stack

import (
	"sync"
)

// Stack is a thread-safe LIFO data structure.
type Stack[T any] struct {
	mu       sync.Mutex
	elements []T
}

// New returns a new empty Stack.
func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

// NewWithCapacity returns a new empty Stack with pre-allocated capacity.
func NewWithCapacity[T any](capacity int) *Stack[T] {
	return &Stack[T]{
		elements: make([]T, 0, capacity),
	}
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(v T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.elements = append(s.elements, v)
}

// Pop removes and returns the top element of the stack.
// Returns (zero-value, false) if the stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	l := len(s.elements)
	if l == 0 {
		var zero T
		return zero, false
	}

	index := l - 1
	v := s.elements[index]

	// Zero out the element to prevent memory leaks (GC can reclaim it)
	var zero T
	s.elements[index] = zero
	s.elements = s.elements[:index]

	return v, true
}

// Peek returns the top element of the stack without removing it.
// Returns (zero-value, false) if the stack is empty.
func (s *Stack[T]) Peek() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	l := len(s.elements)
	if l == 0 {
		var zero T
		return zero, false
	}

	return s.elements[l-1], true
}

// IsEmpty returns true if the stack has no elements.
func (s *Stack[T]) IsEmpty() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.elements) == 0
}

// Len returns the current number of elements in the stack.
func (s *Stack[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.elements)
}

// Clear discards all elements from the stack.
func (s *Stack[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Zero out all elements to assist GC
	var zero T
	for i := range s.elements {
		s.elements[i] = zero
	}
	s.elements = s.elements[:0]
}

// Swap swaps the top two elements of the stack.
// Returns false if the stack has fewer than two elements.
func (s *Stack[T]) Swap() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	l := len(s.elements)
	if l < 2 {
		return false
	}

	s.elements[l-1], s.elements[l-2] = s.elements[l-2], s.elements[l-1]
	return true
}
