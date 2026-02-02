package set

import (
	"sort"
	"testing"
)

type User struct {
	ID   int
	Name string
}

func userEquals(a, b User) bool {
	return a.ID == b.ID
}

func TestAnySet(t *testing.T) {
	s := NewAny(userEquals)
	u1 := User{1, "Alice"}
	u2 := User{2, "Bob"}
	u3 := User{1, "Alice Redux"} // Same ID, should be treated as same item

	// Test Add & Len
	s.Add(u1, u2, u3)
	if s.Len() != 2 {
		t.Errorf("Expected length 2, got %d", s.Len())
	}

	// Test Contains
	if !s.Contains(u1) || !s.Contains(u2) {
		t.Error("Set should contain u1 and u2")
	}
	if !s.Contains(User{1, "Different Name"}) {
		t.Error("Set should contain user with ID 1 regardless of name")
	}

	// Test Remove
	s.Remove(User{2, "Whatever"})
	if s.Len() != 1 {
		t.Errorf("Expected length 1 after remove, got %d", s.Len())
	}
	if s.Contains(u2) {
		t.Error("Set should not contain u2 after removal")
	}

	// Test Pop
	val, ok := s.Pop()
	if !ok {
		t.Error("Pop should return true for non-empty set")
	}
	if s.Len() != 0 {
		t.Errorf("Expected length 0 after pop, got %d", s.Len())
	}
	if s.Contains(val) {
		t.Error("Popped value should no longer be in set")
	}
}

func TestAnySetOperations(t *testing.T) {
	equals := func(a, b int) bool { return a == b }
	s1 := NewAny(equals, 1, 2, 3)
	s2 := NewAny(equals, 3, 4, 5)

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

func TestAnySetComparison(t *testing.T) {
	equals := func(a, b int) bool { return a == b }
	s1 := NewAny(equals, 1, 2)
	s2 := NewAny(equals, 1, 2, 3)

	if !s1.IsSubset(s2) {
		t.Error("s1 should be subset of s2")
	}
	if s2.IsSubset(s1) {
		t.Error("s2 should not be subset of s1")
	}

	s3 := NewAny(equals, 1, 2)
	if !s1.Equal(s3) {
		t.Error("s1 should equal s3")
	}
}

func TestAnySetToSlice(t *testing.T) {
	equals := func(a, b int) bool { return a == b }
	s := NewAny(equals, 1, 2, 3)
	slice := s.ToSlice()
	sort.Ints(slice)
	if len(slice) != 3 || slice[0] != 1 || slice[1] != 2 || slice[2] != 3 {
		t.Errorf("ToSlice returned unexpected result: %v", slice)
	}
}
