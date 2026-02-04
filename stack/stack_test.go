package stack

import (
	"sync"
	"testing"
)

func TestStack(t *testing.T) {
	s := New[int]()

	if !s.IsEmpty() {
		t.Errorf("Expected empty stack")
	}

	s.Push(1)
	s.Push(2)

	if s.Len() != 2 {
		t.Errorf("Expected size 2, got %d", s.Len())
	}

	v, ok := s.Peek()
	if !ok || v != 2 {
		t.Errorf("Expected 2, got %v (ok: %v)", v, ok)
	}

	if !s.Swap() {
		t.Errorf("Expected swap to succeed")
	}

	v, ok = s.Peek()
	if !ok || v != 1 {
		t.Errorf("Expected 1 after swap, got %v (ok: %v)", v, ok)
	}

	v, ok = s.Pop()
	if !ok || v != 1 {
		t.Errorf("Expected 1, got %v (ok: %v)", v, ok)
	}

	v, ok = s.Pop()
	if !ok || v != 2 {
		t.Errorf("Expected 2, got %v (ok: %v)", v, ok)
	}

	if !s.IsEmpty() {
		t.Errorf("Expected empty stack after pops")
	}

	_, ok = s.Pop()
	if ok {
		t.Errorf("Expected ok=false when popping from empty stack")
	}

	s.Push(10)
	s.Clear()
	if s.Len() != 0 {
		t.Errorf("Expected Len 0 after Clear, got %d", s.Len())
	}
}

func TestStackConcurrency(t *testing.T) {
	s := New[int]()
	var wg sync.WaitGroup
	numGoroutines := 100
	numOps := 1000

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				s.Push(j)
			}
		}(i)
	}
	wg.Wait()

	expectedLen := numGoroutines * numOps
	if s.Len() != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, s.Len())
	}

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				s.Pop()
			}
		}()
	}
	wg.Wait()

	if s.Len() != 0 {
		t.Errorf("Expected length 0 after popping all, got %d", s.Len())
	}
}
