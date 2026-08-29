package internal

// ValueMarshaler is able to marshal into a valid simc string
type ValueMarshaler interface{ MarshalSimcValue() (string, error) }

// ValueUnmarshaler is able to unmarshal self from a simc string
type ValueUnmarshaler interface{ UnmarshalSimcValue(string) error }
