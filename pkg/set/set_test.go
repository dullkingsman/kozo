package set

import (
	"sort"
	"testing"
)

func TestSet(t *testing.T) {
	s := New[int]()

	// Test Add & Len
	s.Add(1, 2, 3, 2)
	if s.Len() != 3 {
		t.Errorf("Expected length 3, got %d", s.Len())
	}

	// Test Contains
	if !s.Contains(1) || !s.Contains(2) || !s.Contains(3) {
		t.Error("Set should contain 1, 2, 3")
	}
	if s.Contains(4) {
		t.Error("Set should not contain 4")
	}

	// Test Remove
	s.Remove(2, 4)
	if s.Len() != 2 {
		t.Errorf("Expected length 2 after remove, got %d", s.Len())
	}
	if s.Contains(2) {
		t.Error("Set should not contain 2 after removal")
	}

	// Test Pop
	val, ok := s.Pop()
	if !ok {
		t.Error("Pop should return true for non-empty set")
	}
	if s.Len() != 1 {
		t.Errorf("Expected length 1 after pop, got %d", s.Len())
	}
	if s.Contains(val) {
		t.Error("Popped value should no longer be in set")
	}

	// Test Clear
	s.Clear()
	if !s.IsEmpty() {
		t.Error("Set should be empty after Clear")
	}
}

func TestSetOperations(t *testing.T) {
	s1 := New(1, 2, 3)
	s2 := New(3, 4, 5)

	// Union
	union := s1.Union(s2)
	if union.Len() != 5 {
		t.Errorf("Union should have 5 items, got %d", union.Len())
	}

	// Intersect
	intersect := s1.Intersect(s2)
	if intersect.Len() != 1 || !intersect.Contains(3) {
		t.Error("Intersection should only contain 3")
	}

	// Difference
	diff := s1.Difference(s2)
	if diff.Len() != 2 || !diff.Contains(1) || !diff.Contains(2) {
		t.Error("Difference should contain 1 and 2")
	}

	// SymmetricDifference
	symDiff := s1.SymmetricDifference(s2)
	if symDiff.Len() != 4 || symDiff.Contains(3) {
		t.Error("SymmetricDifference should contain 1, 2, 4, 5 and NOT 3")
	}
}

func TestSetComparison(t *testing.T) {
	s1 := New(1, 2)
	s2 := New(1, 2, 3)

	if !s1.IsSubset(s2) {
		t.Error("s1 should be subset of s2")
	}
	if s2.IsSubset(s1) {
		t.Error("s2 should not be subset of s1")
	}
	if !s2.IsSuperset(s1) {
		t.Error("s2 should be superset of s1")
	}

	s3 := New(1, 2)
	if !s1.Equal(s3) {
		t.Error("s1 should equal s3")
	}
	if s1.Equal(s2) {
		t.Error("s1 should not equal s2")
	}
}

func TestSetIter(t *testing.T) {
	s := New(1, 2, 3)
	sum := 0
	s.Iter(func(v int) bool {
		sum += v
		return true
	})
	if sum != 6 {
		t.Errorf("Expected sum 6, got %d", sum)
	}

	count := 0
	s.Iter(func(v int) bool {
		count++
		return count < 2 // Stop after 2
	})
	if count != 2 {
		t.Errorf("Expected iteration to stop after 2, got %d", count)
	}
}

func TestSetToSlice(t *testing.T) {
	s := New(1, 2, 3)
	slice := s.ToSlice()
	sort.Ints(slice)
	if len(slice) != 3 || slice[0] != 1 || slice[1] != 2 || slice[2] != 3 {
		t.Errorf("ToSlice returned unexpected result: %v", slice)
	}
}
