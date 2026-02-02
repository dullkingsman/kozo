package _range

import (
	"encoding/json"
	"testing"
)

func TestRange_Contains(t *testing.T) {
	less := func(a, b int) bool { return a < b }

	tests := []struct {
		name     string
		r        Range[int]
		val      int
		expected bool
	}{
		{"Closed [10, 20] - inside", Closed(10, 20), 15, true},
		{"Closed [10, 20] - at min", Closed(10, 20), 10, true},
		{"Closed [10, 20] - at max", Closed(10, 20), 20, true},
		{"Closed [10, 20] - below", Closed(10, 20), 9, false},
		{"Closed [10, 20] - above", Closed(10, 20), 21, false},

		{"Open (10, 20) - inside", Open(10, 20), 15, true},
		{"Open (10, 20) - at min", Open(10, 20), 10, false},
		{"Open (10, 20) - at max", Open(10, 20), 20, false},

		{"HalfOpen [10, 20) - at min", HalfOpen(10, 20), 10, true},
		{"HalfOpen [10, 20) - at max", HalfOpen(10, 20), 20, false},

		{"GreaterThan (10) - inside", GreaterThan(10), 11, true},
		{"GreaterThan (10) - at min", GreaterThan(10), 10, false},
		{"AtLeast [10] - at min", AtLeast(10), 10, true},

		{"LessThan (20) - inside", LessThan(20), 19, true},
		{"LessThan (20) - at max", LessThan(20), 20, false},
		{"AtMost [20] - at max", AtMost(20), 20, true},

		{"Any - matches everything", Range[int]{}, 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Contains(tt.val, less); got != tt.expected {
				t.Errorf("Contains(%v) = %v, want %v", tt.val, got, tt.expected)
			}
		})
	}
}

func TestContainsOrdered(t *testing.T) {
	r := Closed(10, 20)
	if !ContainsOrdered(r, 15) {
		t.Error("Expected 15 to be in [10, 20]")
	}
	if ContainsOrdered(r, 25) {
		t.Error("Expected 25 to be outside [10, 20]")
	}
}

func TestRange_Metadata(t *testing.T) {
	r1 := Closed(10, 20)
	if !r1.IsBounded() {
		t.Error("Closed range should be bounded")
	}
	if r1.IsAny() {
		t.Error("Closed range should not be Any")
	}

	r2 := GreaterThan(10)
	if r2.IsBounded() {
		t.Error("Unbounded range should not be bounded")
	}

	r3 := Range[int]{}
	if !r3.IsAny() {
		t.Error("Empty Range should be Any")
	}
}

func TestRange_JSON(t *testing.T) {
	r := Closed(10, 20)
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var r2 Range[int]
	if err := json.Unmarshal(data, &r2); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if !ContainsOrdered(r2, 15) || !ContainsOrdered(r2, 10) || !ContainsOrdered(r2, 20) {
		t.Error("Unmarshaled range does not match expected behavior")
	}
}
