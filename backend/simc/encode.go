package simc

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Marshaler converts a struct to simc format.
type Marshaler interface {
	MarshalSimC() ([]byte, error)
}

func Marshal(input *Input) ([]byte, error) {
	eqp := ""
	chr, err := input.Character.MarshalSimC()
	if err != nil {
		return nil, err
	}
	eqp += string(chr) + "\n"
	for _, eq := range input.Equipment {
		eqStr, err := eq.MarshalSimC()
		if err != nil {
			return nil, err
		}
		eqp += string(eqStr) + "\n"
	}
	opts, err := input.Options.MarshalSimC()
	if err != nil {
		return nil, err
	}
	eqp += string(opts) + "\n"
	return []byte(strings.TrimSpace(eqp)), nil
}

// marshalField handles marshaling a single field based on its type and tag.
func marshalField(field reflect.Value, tag string) ([]byte, error) {
	if !field.IsValid() || field.IsZero() {
		return nil, nil // Skip nil/zero fields
	}
	switch field.Kind() {
	case reflect.String:
		if tag == "" {
			return []byte(field.String()), nil
		}
		return []byte(fmt.Sprintf("%s=%s", tag, field.String())), nil
	case reflect.Int, reflect.Int64:
		if tag == "" {
			return []byte(fmt.Sprintf("%d", field.Int())), nil
		}
		return []byte(fmt.Sprintf("%s=%d", tag, field.Int())), nil
	case reflect.Ptr:
		if field.IsNil() {
			return nil, nil
		}
		return marshalField(field.Elem(), tag)
	case reflect.Slice:
		if field.IsNil() || field.Len() == 0 {
			return nil, nil
		}
		var vals []string
		for i := 0; i < field.Len(); i++ {
			val, err := marshalField(field.Index(i), "")
			if err != nil {
				return nil, err
			}
			vals = append(vals, string(val)) // Remove tag for slices
		}
		return []byte(fmt.Sprintf("%s=%s", tag, strings.Join(vals, "/"))), nil // Use / for multi-values
	case reflect.Struct:
		if marshaler, ok := field.Interface().(Marshaler); ok {
			return marshaler.MarshalSimC()
		}
		// Fallback to generic reflection (not implemented here for brevity)
		return nil, errors.New("nested struct without MarshalSimC")
	default:
		return nil, fmt.Errorf("unsupported type for field %s", tag)
	}
}
