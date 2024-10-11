package jof

import "encoding/json"

// Field is a generic type that can be used to represent a JSON field that may or may not be defined.
// It is useful when you want to distinguish between a field that is defined as null and a field that is not defined at all.
type Field[T any] struct {
	Defined bool
	Value   T
}

func (f *Field[T]) UnmarshalJSON(data []byte) error {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	f.Defined = true
	f.Value = v
	return nil
}

func (f Field[T]) MarshalJSON() ([]byte, error) {
	if !f.Defined {
		return []byte("null"), nil
	}

	return json.Marshal(f.Value)
}

func NewField[T any](v ...T) Field[T] {
	if len(v) == 0 {
		return Field[T]{}
	}

	return Field[T]{
		Defined: true,
		Value:   v[0],
	}
}
