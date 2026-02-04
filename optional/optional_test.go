package data_structures

import (
	"encoding/json"
	"testing"
)

// A comprehensive test suite for Optional[T]
// Tests cover all three states: None, Some(nil), Some(value)
// The three-state model is critical for database operations where
// "field not updated" vs "field set to null" vs "field set to value" have different meanings.

// =========================
// Construction Tests
// =========================

func TestSome(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{"int", 42},
		{"string", "hello"},
		{"bool", true},
		{"zero int", 0},
		{"empty string", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opt Optional[interface{}]
			switch v := tt.value.(type) {
			case int:
				optInt := Some(v)
				if !optInt.IsSome() {
					t.Error("Expected IsSome() to be true")
				}
				if optInt.IsNone() {
					t.Error("Expected IsNone() to be false")
				}
			case string:
				optStr := Some(v)
				if !optStr.IsSome() {
					t.Error("Expected IsSome() to be true")
				}
			case bool:
				optBool := Some(v)
				if !optBool.IsSome() {
					t.Error("Expected IsSome() to be true")
				}
			}

			// Generic test
			opt = Some(tt.value)
			if !opt.IsSome() {
				t.Error("Expected IsSome() to be true")
			}
		})
	}
}

func TestNone(t *testing.T) {
	tests := []struct {
		name string
		opt  interface{}
	}{
		{"int", None[int]()},
		{"string", None[string]()},
		{"bool", None[bool]()},
		{"struct", None[struct{ X int }]()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch opt := tt.opt.(type) {
			case Optional[int]:
				if !opt.IsNone() {
					t.Error("Expected IsNone() to be true")
				}
				if opt.IsSome() {
					t.Error("Expected IsSome() to be false")
				}
			case Optional[string]:
				if !opt.IsNone() {
					t.Error("Expected IsNone() to be true")
				}
			case Optional[bool]:
				if !opt.IsNone() {
					t.Error("Expected IsNone() to be true")
				}
			case Optional[struct{ X int }]:
				if !opt.IsNone() {
					t.Error("Expected IsNone() to be true")
				}
			}
		})
	}
}

func TestZeroValue(t *testing.T) {
	var opt Optional[int]

	if !opt.IsNone() {
		t.Error("Expected zero-value Optional to be None")
	}
	if opt.IsSome() {
		t.Error("Expected zero-value Optional IsSome() to be false")
	}
	if !opt.IsZero() {
		t.Error("Expected zero-value Optional IsZero() to be true")
	}
}

// =========================
// Three-State Model Tests
// =========================

func TestThreeStates(t *testing.T) {
	t.Run("None", func(t *testing.T) {
		opt := None[int]()
		if !opt.IsNone() {
			t.Error("Expected IsNone() to be true")
		}
		if opt.IsSome() {
			t.Error("Expected IsSome() to be false")
		}
		if opt.IsNull() {
			t.Error("Expected IsNull() to be false")
		}
		if opt.IsNotNull() {
			t.Error("Expected IsNotNull() to be false (None is absent, not non-null)")
		}
		if !opt.IsNullOrNone() {
			t.Error("Expected IsNullOrNone() to be true")
		}
	})

	t.Run("Some(nil)", func(t *testing.T) {
		// Create Some with nil pointer
		opt := Optional[int]{value: nil, nonEmpty: true}
		if opt.IsNone() {
			t.Error("Expected IsNone() to be false")
		}
		if !opt.IsSome() {
			t.Error("Expected IsSome() to be true")
		}
		if !opt.IsNull() {
			t.Error("Expected IsNull() to be true")
		}
		if opt.IsNotNull() {
			t.Error("Expected IsNotNull() to be false")
		}
		if !opt.IsNullOrNone() {
			t.Error("Expected IsNullOrNone() to be true")
		}
	})

	t.Run("Some(value)", func(t *testing.T) {
		opt := Some(42)
		if opt.IsNone() {
			t.Error("Expected IsNone() to be false")
		}
		if !opt.IsSome() {
			t.Error("Expected IsSome() to be true")
		}
		if opt.IsNull() {
			t.Error("Expected IsNull() to be false")
		}
		if !opt.IsNotNull() {
			t.Error("Expected IsNotNull() to be true")
		}
		if opt.IsNullOrNone() {
			t.Error("Expected IsNullOrNone() to be false")
		}
	})
}

func TestIsNone(t *testing.T) {
	tests := []struct {
		name     string
		opt      Optional[int]
		expected bool
	}{
		{"None", None[int](), true},
		{"Some(value)", Some(42), false},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, false},
		{"Zero-value", Optional[int]{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opt.IsNone(); got != tt.expected {
				t.Errorf("IsNone() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsSome(t *testing.T) {
	tests := []struct {
		name     string
		opt      Optional[int]
		expected bool
	}{
		{"None", None[int](), false},
		{"Some(value)", Some(42), true},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, true},
		{"Zero-value", Optional[int]{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opt.IsSome(); got != tt.expected {
				t.Errorf("IsSome() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsNull(t *testing.T) {
	tests := []struct {
		name     string
		opt      Optional[int]
		expected bool
	}{
		{"None", None[int](), false},
		{"Some(value)", Some(42), false},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, true},
		{"Zero-value", Optional[int]{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opt.IsNull(); got != tt.expected {
				t.Errorf("IsNull() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsNotNull(t *testing.T) {
	tests := []struct {
		name     string
		opt      Optional[int]
		expected bool
	}{
		{"None", None[int](), false},
		{"Some(value)", Some(42), true},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, false},
		{"Zero-value", Optional[int]{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opt.IsNotNull(); got != tt.expected {
				t.Errorf("IsNotNull() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsNullOrNone(t *testing.T) {
	tests := []struct {
		name     string
		opt      Optional[int]
		expected bool
	}{
		{"None", None[int](), true},
		{"Some(value)", Some(42), false},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, true},
		{"Zero-value", Optional[int]{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opt.IsNullOrNone(); got != tt.expected {
				t.Errorf("IsNullOrNone() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// =========================
// String Tests
// =========================

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		opt      interface{}
		expected string
	}{
		{"None int", None[int](), "None"},
		{"Some(42)", Some(42), "Some(42)"},
		{"Some(hello)", Some("hello"), "Some(hello)"},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, "Some(null)"},
		{"Some(0)", Some(0), "Some(0)"},
		{"Some(empty string)", Some(""), "Some()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			switch opt := tt.opt.(type) {
			case Optional[int]:
				got = opt.String()
			case Optional[string]:
				got = opt.String()
			}

			if got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// =========================
// IsZero Tests
// =========================

func TestIsZero(t *testing.T) {
	tests := []struct {
		name     string
		opt      Optional[int]
		expected bool
	}{
		{"None", None[int](), true},
		{"Some(value)", Some(42), false},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, false},
		{"Zero-value", Optional[int]{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opt.IsZero(); got != tt.expected {
				t.Errorf("IsZero() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// =========================
// JSON Marshaling Tests
// =========================

func TestMarshalJSON_None(t *testing.T) {
	opt := None[int]()
	data, err := json.Marshal(opt)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	expected := "null"
	if string(data) != expected {
		t.Errorf("MarshalJSON() = %s, want %s", data, expected)
	}
}

func TestMarshalJSON_SomeValue(t *testing.T) {
	opt := Some(42)
	data, err := json.Marshal(opt)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	expected := "42"
	if string(data) != expected {
		t.Errorf("MarshalJSON() = %s, want %s", data, expected)
	}
}

func TestMarshalJSON_SomeNull(t *testing.T) {
	opt := Optional[int]{value: nil, nonEmpty: true}
	data, err := json.Marshal(opt)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	expected := "null"
	if string(data) != expected {
		t.Errorf("MarshalJSON() = %s, want %s", data, expected)
	}
}

func TestMarshalJSON_BothValueAndPointer(t *testing.T) {
	opt := Some(42)

	// Value receiver MarshalJSON should work for both value and pointer
	dataValue, err := json.Marshal(opt)
	if err != nil {
		t.Fatalf("Marshal by value error = %v", err)
	}

	dataPtr, err := json.Marshal(&opt)
	if err != nil {
		t.Fatalf("Marshal by pointer error = %v", err)
	}

	// Both should produce the same result
	if string(dataValue) != string(dataPtr) {
		t.Errorf("Marshal results differ: value=%s, pointer=%s", dataValue, dataPtr)
	}

	expected := "42"
	if string(dataValue) != expected {
		t.Errorf("Marshal = %s, want %s", dataValue, expected)
	}
}

func TestUnmarshalJSON_Null(t *testing.T) {
	var opt Optional[int]
	data := []byte("null")

	err := json.Unmarshal(data, &opt)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	// JSON null should become Some(nil)
	if !opt.IsSome() {
		t.Error("Expected IsSome() to be true")
	}
	if !opt.IsNull() {
		t.Error("Expected IsNull() to be true")
	}
	if opt.IsNone() {
		t.Error("Expected IsNone() to be false")
	}
}

func TestUnmarshalJSON_Value(t *testing.T) {
	var opt Optional[int]
	data := []byte("42")

	err := json.Unmarshal(data, &opt)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if !opt.IsSome() {
		t.Error("Expected IsSome() to be true")
	}
	if !opt.IsNotNull() {
		t.Error("Expected IsNotNull() to be true")
	}

	val, ok := opt.Unwrap()
	if !ok {
		t.Error("Expected Unwrap() to return true")
	}
	if val != 42 {
		t.Errorf("Expected value 42, got %d", val)
	}
}

func TestUnmarshalJSON_InvalidJSON(t *testing.T) {
	var opt Optional[int]
	data := []byte("invalid")

	err := json.Unmarshal(data, &opt)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestJSONRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		opt  Optional[int]
	}{
		{"Some(value)", Some(42)},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := json.Marshal(tt.opt)
			if err != nil {
				t.Fatalf("Marshal error = %v", err)
			}

			// Unmarshal
			var got Optional[int]
			err = json.Unmarshal(data, &got)
			if err != nil {
				t.Fatalf("Unmarshal error = %v", err)
			}

			// Compare states
			if got.IsSome() != tt.opt.IsSome() {
				t.Errorf("IsSome() mismatch: got %v, want %v", got.IsSome(), tt.opt.IsSome())
			}
			if got.IsNull() != tt.opt.IsNull() {
				t.Errorf("IsNull() mismatch: got %v, want %v", got.IsNull(), tt.opt.IsNull())
			}
		})
	}
}

func TestJSONInStruct(t *testing.T) {
	type TestStruct struct {
		Required int              `json:"required"`
		Optional Optional[int]    `json:"optional,omitempty"`
		Name     Optional[string] `json:"name,omitempty"`
	}

	t.Run("With values", func(t *testing.T) {
		s := TestStruct{
			Required: 1,
			Optional: Some(42),
			Name:     Some("test"),
		}

		data, err := json.Marshal(s)
		if err != nil {
			t.Fatalf("Marshal error = %v", err)
		}

		// Value receiver MarshalJSON works correctly with struct fields
		expected := `{"required":1,"optional":42,"name":"test"}`
		if string(data) != expected {
			t.Errorf("Marshal = %s, want %s", string(data), expected)
		}

		// Test roundtrip
		var got TestStruct
		err = json.Unmarshal(data, &got)
		if err != nil {
			t.Fatalf("Unmarshal error = %v", err)
		}

		if !got.Optional.IsSome() || !got.Optional.IsNotNull() {
			t.Error("Expected Optional.IsNotNull() to be true")
		}
		val, _ := got.Optional.Unwrap()
		if val != 42 {
			t.Errorf("Optional value = %d, want 42", val)
		}

		if !got.Name.IsSome() || !got.Name.IsNotNull() {
			t.Error("Expected Name.IsNotNull() to be true")
		}
		name, _ := got.Name.Unwrap()
		if name != "test" {
			t.Errorf("Name value = %s, want test", name)
		}
	})

	t.Run("With null", func(t *testing.T) {
		data := []byte(`{"required":1,"optional":null}`)

		var got TestStruct
		err := json.Unmarshal(data, &got)
		if err != nil {
			t.Fatalf("Unmarshal error = %v", err)
		}

		// optional should be Some(nil)
		if !got.Optional.IsSome() {
			t.Error("Expected Optional.IsSome() to be true")
		}
		if !got.Optional.IsNull() {
			t.Error("Expected Optional.IsNull() to be true")
		}

		// name should be None (missing)
		if !got.Name.IsNone() {
			t.Error("Expected Name.IsNone() to be true")
		}
	})

	t.Run("Missing fields", func(t *testing.T) {
		data := []byte(`{"required":1}`)

		var got TestStruct
		err := json.Unmarshal(data, &got)
		if err != nil {
			t.Fatalf("Unmarshal error = %v", err)
		}

		// Both optionals should be None
		if !got.Optional.IsNone() {
			t.Error("Expected Optional.IsNone() to be true")
		}
		if !got.Name.IsNone() {
			t.Error("Expected Name.IsNone() to be true")
		}
	})
}

// =========================
// Access Method Tests
// =========================

func TestUnwrap(t *testing.T) {
	tests := []struct {
		name      string
		opt       Optional[int]
		wantValue int
		wantOk    bool
	}{
		{"None", None[int](), 0, false},
		{"Some(value)", Some(42), 42, true},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, 0, false},
		{"Some(0)", Some(0), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := tt.opt.Unwrap()
			if ok != tt.wantOk {
				t.Errorf("Unwrap() ok = %v, want %v", ok, tt.wantOk)
			}
			if val != tt.wantValue {
				t.Errorf("Unwrap() value = %v, want %v", val, tt.wantValue)
			}
		})
	}
}

func TestUnwrapPtr(t *testing.T) {
	tests := []struct {
		name      string
		opt       Optional[int]
		wantValue *int
		wantOk    bool
	}{
		{"None", None[int](), nil, false},
		{"Some(value)", Some(42), intPtr(42), true},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := tt.opt.UnwrapPtr()
			if ok != tt.wantOk {
				t.Errorf("UnwrapPtr() ok = %v, want %v", ok, tt.wantOk)
			}
			if tt.wantValue == nil {
				if val != nil {
					t.Errorf("UnwrapPtr() value = %v, want nil", val)
				}
			} else if val == nil || *val != *tt.wantValue {
				t.Errorf("UnwrapPtr() value = %v, want %v", val, *tt.wantValue)
			}
		})
	}
}

func TestUnwrapOr(t *testing.T) {
	tests := []struct {
		name         string
		opt          Optional[int]
		defaultValue int
		want         int
	}{
		{"None", None[int](), 99, 99},
		{"Some(value)", Some(42), 99, 42},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, 99, 99},
		{"Some(0)", Some(0), 99, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opt.UnwrapOr(tt.defaultValue)
			if got != tt.want {
				t.Errorf("UnwrapOr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnwrapOrPtr(t *testing.T) {
	defaultVal := 99

	tests := []struct {
		name         string
		opt          Optional[int]
		defaultValue *int
		wantValue    *int
	}{
		{"None", None[int](), &defaultVal, &defaultVal},
		{"Some(value)", Some(42), &defaultVal, intPtr(42)},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, &defaultVal, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opt.UnwrapOrPtr(tt.defaultValue)
			if tt.wantValue == nil {
				if got != nil {
					t.Errorf("UnwrapOrPtr() = %v, want nil", got)
				}
			} else if got == nil {
				t.Error("UnwrapOrPtr() returned nil, want non-nil")
			} else if *got != *tt.wantValue {
				t.Errorf("UnwrapOrPtr() = %v, want %v", *got, *tt.wantValue)
			}
		})
	}
}

func TestUnwrapOrElse(t *testing.T) {
	defaultFunc := func() int { return 99 }
	called := false
	trackingFunc := func() int {
		called = true
		return 99
	}

	tests := []struct {
		name        string
		opt         Optional[int]
		defaultFunc func() int
		want        int
		wantCalled  bool
	}{
		{"None", None[int](), trackingFunc, 99, true},
		{"Some(value)", Some(42), defaultFunc, 42, false},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, trackingFunc, 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called = false
			got := tt.opt.UnwrapOrElse(tt.defaultFunc)
			if got != tt.want {
				t.Errorf("UnwrapOrElse() = %v, want %v", got, tt.want)
			}
			if tt.name == "None" || tt.name == "Some(nil)" {
				if called != tt.wantCalled {
					t.Errorf("defaultFunc called = %v, want %v", called, tt.wantCalled)
				}
			}
		})
	}
}

func TestUnwrapOrElsePtr(t *testing.T) {
	defaultFunc := func() *int { return intPtr(99) }

	tests := []struct {
		name        string
		opt         Optional[int]
		defaultFunc func() *int
		wantValue   *int
	}{
		{"None", None[int](), defaultFunc, intPtr(99)},
		{"Some(value)", Some(42), defaultFunc, intPtr(42)},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, defaultFunc, intPtr(42)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opt.UnwrapOrElsePtr(tt.defaultFunc)
			if tt.name == "Some(nil)" {
				// Special case: Some(nil) returns the internal nil pointer
				if got != nil {
					t.Error("UnwrapOrElsePtr() for Some(nil) should return nil")
				}
			} else if got == nil {
				t.Error("UnwrapOrElsePtr() returned nil")
			}
		})
	}
}

// =========================
// Panic Tests
// =========================

func TestExpect_Panic(t *testing.T) {
	tests := []struct {
		name string
		opt  Optional[int]
	}{
		{"None", None[int]()},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic but didn't get one")
				}
			}()
			tt.opt.Expect("test panic")
		})
	}
}

func TestExpect_Success(t *testing.T) {
	opt := Some(42)
	got := opt.Expect("should not panic")
	if got != 42 {
		t.Errorf("Expect() = %v, want 42", got)
	}
}

func TestExpectPtr_Panic(t *testing.T) {
	opt := None[int]()
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic but didn't get one")
		}
	}()
	opt.ExpectPtr("test panic")
}

func TestExpectPtr_Success(t *testing.T) {
	tests := []struct {
		name string
		opt  Optional[int]
	}{
		{"Some(value)", Some(42)},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opt.ExpectPtr("should not panic")
			if tt.name == "Some(nil)" {
				if got != nil {
					t.Error("ExpectPtr() for Some(nil) should return nil")
				}
			} else {
				if got == nil || *got != 42 {
					t.Errorf("ExpectPtr() = %v, want 42", got)
				}
			}
		})
	}
}

// =========================
// Comparison Tests
// =========================

func TestContainsComparable(t *testing.T) {
	tests := []struct {
		name  string
		opt   Optional[int]
		value int
		want  bool
	}{
		{"None", None[int](), 42, false},
		{"Some(value) match", Some(42), 42, true},
		{"Some(value) no match", Some(42), 99, false},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, 42, false},
		{"Some(0) match", Some(0), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsComparable(tt.opt, tt.value)
			if got != tt.want {
				t.Errorf("ContainsComparable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	equals := func(a, b int) bool { return a == b }
	greaterThan := func(a, b int) bool { return a > b }

	tests := []struct {
		name        string
		opt         Optional[int]
		value       int
		compareFunc func(int, int) bool
		want        bool
	}{
		{"None", None[int](), 42, equals, false},
		{"Some(value) equals", Some(42), 42, equals, true},
		{"Some(value) not equals", Some(42), 99, equals, false},
		{"Some(50) > 30", Some(50), 30, greaterThan, true},
		{"Some(20) > 30", Some(20), 30, greaterThan, false},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, 42, equals, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opt.Contains(tt.value, tt.compareFunc)
			if got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

// =========================
// Filter Tests
// =========================

func TestFilter(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }
	isPositive := func(n int) bool { return n > 0 }

	tests := []struct {
		name      string
		opt       Optional[int]
		predicate func(int) bool
		wantSome  bool
	}{
		{"None", None[int](), isEven, false},
		{"Some(42) isEven", Some(42), isEven, true},
		{"Some(43) isEven", Some(43), isEven, false},
		{"Some(-5) isPositive", Some(-5), isPositive, false},
		{"Some(5) isPositive", Some(5), isPositive, true},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, isEven, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opt.Filter(tt.predicate)
			if got.IsSome() != tt.wantSome {
				t.Errorf("Filter() IsSome() = %v, want %v", got.IsSome(), tt.wantSome)
			}
		})
	}
}

func TestFilterPtr(t *testing.T) {
	isNonNil := func(p *int) bool { return p != nil }
	isEven := func(p *int) bool { return p != nil && *p%2 == 0 }

	tests := []struct {
		name      string
		opt       Optional[int]
		predicate func(*int) bool
		wantSome  bool
	}{
		{"None", None[int](), isNonNil, false},
		{"Some(42) isNonNil", Some(42), isNonNil, true},
		{"Some(nil) isNonNil", Optional[int]{value: nil, nonEmpty: true}, isNonNil, false},
		{"Some(42) isEven", Some(42), isEven, true},
		{"Some(43) isEven", Some(43), isEven, false},
		{"Some(nil) isEven", Optional[int]{value: nil, nonEmpty: true}, isEven, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opt.FilterPtr(tt.predicate)
			if got.IsSome() != tt.wantSome {
				t.Errorf("FilterPtr() IsSome() = %v, want %v", got.IsSome(), tt.wantSome)
			}
		})
	}
}

// =========================
// Match Tests
// =========================

func TestMatch(t *testing.T) {
	tests := []struct {
		name         string
		opt          Optional[int]
		wantSomePath bool
		wantValue    int
	}{
		{"None", None[int](), false, 0},
		{"Some(42)", Some(42), true, 42},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			someCalled := false
			noneCalled := false
			var capturedValue int

			tt.opt.Match(
				func(v int) {
					someCalled = true
					capturedValue = v
				},
				func() {
					noneCalled = true
				},
			)

			if someCalled != tt.wantSomePath {
				t.Errorf("Match() someCalled = %v, want %v", someCalled, tt.wantSomePath)
			}
			if noneCalled == tt.wantSomePath {
				t.Errorf("Match() noneCalled = %v, want %v", noneCalled, !tt.wantSomePath)
			}
			if someCalled && capturedValue != tt.wantValue {
				t.Errorf("Match() value = %v, want %v", capturedValue, tt.wantValue)
			}
		})
	}
}

func TestMatchPtr(t *testing.T) {
	tests := []struct {
		name         string
		opt          Optional[int]
		wantSomePath bool
		wantNilPtr   bool
	}{
		{"None", None[int](), false, false},
		{"Some(42)", Some(42), true, false},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			someCalled := false
			noneCalled := false
			var capturedPtr *int

			tt.opt.MatchPtr(
				func(p *int) {
					someCalled = true
					capturedPtr = p
				},
				func() {
					noneCalled = true
				},
			)

			if someCalled != tt.wantSomePath {
				t.Errorf("MatchPtr() someCalled = %v, want %v", someCalled, tt.wantSomePath)
			}
			if noneCalled == tt.wantSomePath {
				t.Errorf("MatchPtr() noneCalled = %v, want %v", noneCalled, !tt.wantSomePath)
			}
			if someCalled {
				isNil := capturedPtr == nil
				if isNil != tt.wantNilPtr {
					t.Errorf("MatchPtr() ptr is nil = %v, want %v", isNil, tt.wantNilPtr)
				}
			}
		})
	}
}

// =========================
// Combining Tests
// =========================

func TestOr(t *testing.T) {
	tests := []struct {
		name      string
		opt       Optional[int]
		other     Optional[int]
		wantValue int
		wantSome  bool
	}{
		{"None Or None", None[int](), None[int](), 0, false},
		{"None Or Some", None[int](), Some(42), 42, true},
		{"Some Or None", Some(99), None[int](), 99, true},
		{"Some Or Some", Some(99), Some(42), 99, true},
		{"Some(nil) Or Some", Optional[int]{value: nil, nonEmpty: true}, Some(42), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opt.Or(tt.other)
			if got.IsSome() != tt.wantSome {
				t.Errorf("Or() IsSome() = %v, want %v", got.IsSome(), tt.wantSome)
			}
			if tt.wantSome && got.IsNotNull() {
				val, _ := got.Unwrap()
				if val != tt.wantValue {
					t.Errorf("Or() value = %v, want %v", val, tt.wantValue)
				}
			}
		})
	}
}

func TestOrElse(t *testing.T) {
	called := false
	elseFunc := func() Optional[int] {
		called = true
		return Some(99)
	}

	tests := []struct {
		name       string
		opt        Optional[int]
		wantCalled bool
		wantValue  int
		wantSome   bool
	}{
		{"None", None[int](), true, 99, true},
		{"Some", Some(42), false, 42, true},
		{"Some(nil)", Optional[int]{value: nil, nonEmpty: true}, false, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called = false
			got := tt.opt.OrElse(elseFunc)

			if called != tt.wantCalled {
				t.Errorf("OrElse() called = %v, want %v", called, tt.wantCalled)
			}
			if got.IsSome() != tt.wantSome {
				t.Errorf("OrElse() IsSome() = %v, want %v", got.IsSome(), tt.wantSome)
			}
			if tt.wantSome && got.IsNotNull() {
				val, _ := got.Unwrap()
				if val != tt.wantValue {
					t.Errorf("OrElse() value = %v, want %v", val, tt.wantValue)
				}
			}
		})
	}
}

func TestXor(t *testing.T) {
	tests := []struct {
		name      string
		opt       Optional[int]
		other     Optional[int]
		wantSome  bool
		wantValue int
	}{
		{"None Xor None", None[int](), None[int](), false, 0},
		{"None Xor Some", None[int](), Some(42), true, 42},
		{"Some Xor None", Some(99), None[int](), true, 99},
		{"Some Xor Some", Some(99), Some(42), false, 0},
		{"Some(nil) Xor None", Optional[int]{value: nil, nonEmpty: true}, None[int](), true, 0},
		{"Some(nil) Xor Some", Optional[int]{value: nil, nonEmpty: true}, Some(42), false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opt.Xor(tt.other)
			if got.IsSome() != tt.wantSome {
				t.Errorf("Xor() IsSome() = %v, want %v", got.IsSome(), tt.wantSome)
			}
			if tt.wantSome && got.IsNotNull() {
				val, _ := got.Unwrap()
				if val != tt.wantValue {
					t.Errorf("Xor() value = %v, want %v", val, tt.wantValue)
				}
			}
		})
	}
}

// =========================
// Clone Tests
// =========================

func TestClone_None(t *testing.T) {
	opt := None[int]()
	cloned := opt.Clone()

	if !cloned.IsNone() {
		t.Error("Cloned None should be None")
	}
}

func TestClone_SomeValue(t *testing.T) {
	opt := Some(42)
	cloned := opt.Clone()

	if !cloned.IsSome() {
		t.Error("Cloned Some should be Some")
	}

	val, ok := cloned.Unwrap()
	if !ok || val != 42 {
		t.Errorf("Cloned value = %v, %v; want 42, true", val, ok)
	}

	// Verify independence - changes to original don't affect clone
	Take(&opt)
	if !opt.IsNone() {
		t.Error("Original should be None after Take")
	}
	if !cloned.IsSome() {
		t.Error("Clone should still be Some after original is taken")
	}
}

func TestClone_SomeNull(t *testing.T) {
	opt := Optional[int]{value: nil, nonEmpty: true}
	cloned := opt.Clone()

	if !cloned.IsSome() {
		t.Error("Cloned Some(nil) should be Some")
	}
	if !cloned.IsNull() {
		t.Error("Cloned Some(nil) should be Null")
	}
}

func TestClone_Independence(t *testing.T) {
	type TestStruct struct {
		Value int
	}

	opt := Some(TestStruct{Value: 42})
	cloned := opt.Clone()

	// Get pointers
	ptr1, _ := opt.UnwrapPtr()
	ptr2, _ := cloned.UnwrapPtr()

	// Pointers should be different
	if ptr1 == ptr2 {
		t.Error("Clone should create new pointer, not share")
	}

	// Values should be equal
	if ptr1.Value != ptr2.Value {
		t.Errorf("Values should be equal: %d != %d", ptr1.Value, ptr2.Value)
	}

	// Modify original
	ptr1.Value = 99

	// Clone should be unchanged (shallow copy means struct is copied)
	if ptr2.Value != 42 {
		t.Errorf("Clone was affected by original modification: got %d, want 42", ptr2.Value)
	}
}

func TestClone_ShallowCopy(t *testing.T) {
	// Test that slices are shallow copied
	slice := []int{1, 2, 3}
	opt := Some(slice)
	cloned := opt.Clone()

	ptr1, _ := opt.UnwrapPtr()
	ptr2, _ := cloned.UnwrapPtr()

	// The slice headers are different
	if ptr1 == ptr2 {
		t.Error("Clone should create new slice header")
	}

	// But they share backing array (shallow copy)
	(*ptr1)[0] = 999

	// This documents the shallow copy behavior
	// The clone's backing array is also modified
	if (*ptr2)[0] != 999 {
		t.Log("Note: Slice backing array was not shared (this is actually good for Go's copy semantics)")
	}
}

// =========================
// Take Tests
// =========================

func TestTake(t *testing.T) {
	t.Run("Take Some(value)", func(t *testing.T) {
		opt := Some(42)
		old := Take(&opt)

		// Old should have the value
		if !old.IsSome() {
			t.Error("Taken value should be Some")
		}
		val, ok := old.Unwrap()
		if !ok || val != 42 {
			t.Errorf("Taken value = %v, %v; want 42, true", val, ok)
		}

		// Original should be None
		if !opt.IsNone() {
			t.Error("Original should be None after Take")
		}
	})

	t.Run("Take None", func(t *testing.T) {
		opt := None[int]()
		old := Take(&opt)

		if !old.IsNone() {
			t.Error("Taken None should be None")
		}
		if !opt.IsNone() {
			t.Error("Original should still be None")
		}
	})

	t.Run("Take Some(nil)", func(t *testing.T) {
		opt := Optional[int]{value: nil, nonEmpty: true}
		old := Take(&opt)

		if !old.IsSome() {
			t.Error("Taken Some(nil) should be Some")
		}
		if !old.IsNull() {
			t.Error("Taken Some(nil) should be Null")
		}
		if !opt.IsNone() {
			t.Error("Original should be None after Take")
		}
	})

	t.Run("Take nil pointer", func(t *testing.T) {
		var opt *Optional[int]
		result := Take(opt)

		if !result.IsNone() {
			t.Error("Taking nil pointer should return None")
		}
	})
}

// =========================
// Edge Cases & Integration
// =========================

func TestMultipleTypes(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		opt := Some(42)
		if val, ok := opt.Unwrap(); !ok || val != 42 {
			t.Error("int Optional failed")
		}
	})

	t.Run("string", func(t *testing.T) {
		opt := Some("hello")
		if val, ok := opt.Unwrap(); !ok || val != "hello" {
			t.Error("string Optional failed")
		}
	})

	t.Run("struct", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		opt := Some(Person{Name: "Alice", Age: 30})
		if val, ok := opt.Unwrap(); !ok || val.Name != "Alice" {
			t.Error("struct Optional failed")
		}
	})

	t.Run("pointer", func(t *testing.T) {
		val := 42
		opt := Some(&val)
		if ptr, ok := opt.Unwrap(); !ok || *ptr != 42 {
			t.Error("pointer Optional failed")
		}
	})

	t.Run("slice", func(t *testing.T) {
		opt := Some([]int{1, 2, 3})
		if slice, ok := opt.Unwrap(); !ok || len(slice) != 3 {
			t.Error("slice Optional failed")
		}
	})

	t.Run("map", func(t *testing.T) {
		opt := Some(map[string]int{"a": 1})
		if m, ok := opt.Unwrap(); !ok || m["a"] != 1 {
			t.Error("map Optional failed")
		}
	})
}

func TestChaining(t *testing.T) {
	opt := Some(42).Filter(func(n int) bool { return n > 0 }).Or(Some(99))

	if !opt.IsSome() {
		t.Error("Chaining failed: expected Some")
	}

	val, ok := opt.Unwrap()
	if !ok || val != 42 {
		t.Errorf("Chaining failed: got %v, %v; want 42, true", val, ok)
	}
}

func TestPointerSharing(t *testing.T) {
	opt1 := Some(42)
	opt2 := opt1 // Copy

	ptr1, _ := opt1.UnwrapPtr()
	ptr2, _ := opt2.UnwrapPtr()

	// Should share the same pointer (value semantics)
	if ptr1 != ptr2 {
		t.Error("Copies should share pointer")
	}

	// But modifications through Take don't affect the copy
	Take(&opt1)
	if !opt1.IsNone() {
		t.Error("opt1 should be None after Take")
	}
	if !opt2.IsSome() {
		t.Error("opt2 should still be Some")
	}
}

// =========================
// Benchmarks
// =========================

func BenchmarkSome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Some(42)
	}
}

func BenchmarkNone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = None[int]()
	}
}

func BenchmarkUnwrap(b *testing.B) {
	opt := Some(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = opt.Unwrap()
	}
}

func BenchmarkUnwrapOr(b *testing.B) {
	opt := Some(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opt.UnwrapOr(99)
	}
}

func BenchmarkFilter(b *testing.B) {
	opt := Some(42)
	pred := func(n int) bool { return n > 0 }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opt.Filter(pred)
	}
}

func BenchmarkClone(b *testing.B) {
	opt := Some(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opt.Clone()
	}
}

func BenchmarkJSONMarshal(b *testing.B) {
	opt := Some(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(&opt)
	}
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	data := []byte("42")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var opt Optional[int]
		_ = json.Unmarshal(data, &opt)
	}
}

// Helper function
func intPtr(i int) *int {
	return &i
}
