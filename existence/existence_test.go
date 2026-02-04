package existence

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestExistenceClaim_In(t *testing.T) {
	ec := In(1, 2, 3)
	if !ec.Contains {
		t.Error("Expected Contains to be true")
	}
	if len(ec.Values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(ec.Values))
	}
}

func TestExistenceClaim_NotIn(t *testing.T) {
	ec := NotIn(1, 2, 3)
	if ec.Contains {
		t.Error("Expected Contains to be false")
	}
	if len(ec.Values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(ec.Values))
	}
}

func TestExistenceClaim_Check(t *testing.T) {
	equals := func(a, b int) bool { return a == b }

	t.Run("In matches", func(t *testing.T) {
		ec := In(1, 2, 3)
		if !ec.Check(2, equals) {
			t.Error("Expected Check(2) to be true for In(1,2,3)")
		}
	})

	t.Run("In does not match", func(t *testing.T) {
		ec := In(1, 2, 3)
		if ec.Check(4, equals) {
			t.Error("Expected Check(4) to be false for In(1,2,3)")
		}
	})

	t.Run("NotIn matches", func(t *testing.T) {
		ec := NotIn(1, 2, 3)
		if !ec.Check(4, equals) {
			t.Error("Expected Check(4) to be true for NotIn(1,2,3)")
		}
	})

	t.Run("NotIn does not match", func(t *testing.T) {
		ec := NotIn(1, 2, 3)
		if ec.Check(2, equals) {
			t.Error("Expected Check(2) to be false for NotIn(1,2,3)")
		}
	})
}

func TestExistenceClaim_CheckComparable(t *testing.T) {
	t.Run("In matches", func(t *testing.T) {
		ec := In("a", "b")
		if !CheckComparable(ec, "a") {
			t.Error("Expected CheckComparable to be true")
		}
	})

	t.Run("NotIn matches", func(t *testing.T) {
		ec := NotIn("a", "b")
		if !CheckComparable(ec, "c") {
			t.Error("Expected CheckComparable to be true")
		}
	})
}

func TestExistenceClaim_IsEmpty(t *testing.T) {
	ec := In[int]()
	if !ec.IsEmpty() {
		t.Error("Expected IsEmpty to be true")
	}

	ec2 := In(1)
	if ec2.IsEmpty() {
		t.Error("Expected IsEmpty to be false")
	}
}

func TestExistenceClaim_Len(t *testing.T) {
	ec := In(1, 2, 3)
	if ec.Len() != 3 {
		t.Errorf("Expected Len to be 3, got %d", ec.Len())
	}
}

func TestExistenceClaim_Negate(t *testing.T) {
	ec := In(1, 2)
	neg := ec.Negate()
	if neg.Contains {
		t.Error("Expected Negate() to have Contains=false")
	}
	if len(neg.Values) != 2 {
		t.Error("Values should be preserved after Negate()")
	}
}

func TestExistenceClaim_Apply(t *testing.T) {
	ec := In(1, 2)
	equals := func(a, b int) bool { return a == b }
	input := []int{1, 2, 3, 4, 1}
	expected := []int{1, 2, 1}

	result := ec.Apply(input, equals)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Apply failed. Got %v, want %v", result, expected)
	}
}

func TestExistenceClaim_JSON(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		ec := In(1, 2)
		data, err := json.Marshal(ec)
		if err != nil {
			t.Fatal(err)
		}
		expected := `{"in":[1,2],"contains":true}`
		if string(data) != expected {
			t.Errorf("Marshal mismatch. Got %s, want %s", string(data), expected)
		}
	})

	t.Run("Unmarshal", func(t *testing.T) {
		input := `{"in":[1,2],"contains":false}`
		var ec ExistenceClaim[int]
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatal(err)
		}
		if ec.Contains {
			t.Error("Expected Contains to be false")
		}
		if len(ec.Values) != 2 || ec.Values[0] != 1 || ec.Values[1] != 2 {
			t.Errorf("Values mismatch: %v", ec.Values)
		}
	})
}
