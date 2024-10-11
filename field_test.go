package jof

import (
	"encoding/json"
	"testing"
)

func TestField_UnmarshalJSON(t *testing.T) {
	type User struct {
		ID       int            `json:"id"`
		Name     Field[string]  `json:"name"`
		Username Field[*string] `json:"username"`
		Age      Field[int]     `json:"age"`
	}
	t.Run("UndefinedFields", func(t *testing.T) {
		j := []byte(`{"id": 1}`)
		var u User
		if err := json.Unmarshal(j, &u); err != nil {
			t.Fatal(err)
		}

		if u.Name.Defined {
			t.Error("Expected Name to be undefined")
		}
		if u.Username.Defined {
			t.Error("Expected Username to be undefined")
		}
	})

	t.Run("DefinedField_Null", func(t *testing.T) {
		j := []byte(`{"id": 1, "username": null}`)
		var u User
		if err := json.Unmarshal(j, &u); err != nil {
			t.Fatal(err)
		}

		if !u.Username.Defined {
			t.Error("Expected Username to be defined")
		}

		if u.Username.Value != nil {
			t.Errorf("Expected Username to be nil, got %v", u.Username.Value)
		}
	})

	t.Run("DefinedField_NotNull", func(t *testing.T) {
		j := []byte(`{"id": 1, "username": "alice"}`)
		var u User
		if err := json.Unmarshal(j, &u); err != nil {
			t.Fatal(err)
		}

		if !u.Username.Defined {
			t.Error("Expected Username to be defined")
		}

		if u.Username.Value == nil {
			t.Error("Expected Username to be non-nil")
		}

		if *u.Username.Value != "alice" {
			t.Errorf("Expected Username to be alice, got %v", *u.Username.Value)
		}
	})

	t.Run("DefinedField_Invalid_Int", func(t *testing.T) {
		j := []byte(`{"id": 1, "name": "Alice"}`)
		var u User
		if err := json.Unmarshal(j, &u); err != nil {
			t.Fatal(err)
		}

		if !u.Name.Defined {
			t.Error("Expected Name to be defined")
		}

		if u.Name.Value != "Alice" {
			t.Errorf("Expected Name to be Alice, got %v", u.Name.Value)
		}
	})

	t.Run("DefinedField_Valid_Int", func(t *testing.T) {
		j := []byte(`{"id": 1, "age": "Alice"}`)
		var u User
		if err := json.Unmarshal(j, &u); err == nil {
			t.Error("Expected error")
		}
	})
}

func TestField_MarshalJSON(t *testing.T) {
	t.Run("MarshalUndefinedField", func(t *testing.T) {
		field := NewField[int]() // Undefined field
		data, err := json.Marshal(field)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		expected := "null"
		if string(data) != expected {
			t.Errorf("Expected %v, got %v", expected, string(data))
		}
	})
	t.Run("MarshalDefinedField_Null", func(t *testing.T) {
		field := NewField[*int](nil) // Defined field with nil value
		data, err := json.Marshal(field)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		expected := "null"
		if string(data) != expected {
			t.Errorf("Expected %v, got %v", expected, string(data))
		}
	})

	t.Run("MarshalDefinedField_Integer", func(t *testing.T) {
		field := NewField(42) // Defined field with integer value
		data, err := json.Marshal(field)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		expected := "42"
		if string(data) != expected {
			t.Errorf("Expected %v, got %v", expected, string(data))
		}
	})

	t.Run("MarshalDefinedField_String", func(t *testing.T) {
		field := NewField("hello") // Defined field with string value
		data, err := json.Marshal(field)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		expected := `"hello"`
		if string(data) != expected {
			t.Errorf("Expected %v, got %v", expected, string(data))
		}
	})

	t.Run("MarshalDefinedField_Boolean", func(t *testing.T) {
		field := NewField(true) // Defined field with boolean value
		data, err := json.Marshal(field)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		expected := "true"
		if string(data) != expected {
			t.Errorf("Expected %v, got %v", expected, string(data))
		}
	})
}

func TestNewField(t *testing.T) {
	t.Run("WithoutArgument", func(t *testing.T) {
		field := NewField[int]()
		if field.Defined {
			t.Errorf("Expected Defined to be false, got true")
		}
	})

	t.Run("WithArgument", func(t *testing.T) {
		field := NewField(42)
		if !field.Defined {
			t.Errorf("Expected Defined to be true, got false")
		}
		if field.Value != 42 {
			t.Errorf("Expected Value to be 42, got %v", field.Value)
		}
	})

	t.Run("WithDifferentTypes_String", func(t *testing.T) {
		stringField := NewField("hello")
		if !stringField.Defined {
			t.Errorf("Expected Defined to be true for string field, got false")
		}
		if stringField.Value != "hello" {
			t.Errorf("Expected Value to be 'hello', got %v", stringField.Value)
		}
	})

	t.Run("WithDifferentTypes_Float", func(t *testing.T) {
		floatField := NewField(3.14)
		if !floatField.Defined {
			t.Errorf("Expected Defined to be true for float field, got false")
		}
		if floatField.Value != 3.14 {
			t.Errorf("Expected Value to be 3.14, got %v", floatField.Value)
		}
	})
}
