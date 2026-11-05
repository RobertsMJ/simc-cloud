package internal

import "errors"

var (
	ErrUnmarshalIntoNilPtr     = errors.New("attempting to unmarshal into nil pointer")
	ErrNotPointer              = errors.New("value must be a pointer")
	ErrInvalidStatementForType = errors.New("invalid statement for type")
	ErrInvalidKeyForType       = errors.New("invalid statement key for type")
)
