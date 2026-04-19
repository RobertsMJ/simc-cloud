package simc

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Unmarshaler parses simc format into a struct.
type Unmarshaler interface {
	UnmarshalSimC(data []byte) error
}

// Unmarshal parses the full simc string
func Unmarshal(data []byte, input *Input) error {
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if IsCharacterValue(line) { // Character line
			if err := input.Character.UnmarshalSimC([]byte(line)); err != nil {
				return err
			}
		} else if IsEquipmentValue(line) { // Equipment line
			if input.Equipment == nil {
				input.Equipment = make(map[EquipmentSlot]Equipment)
			}
			var eq Equipment
			if err := eq.UnmarshalSimC([]byte(line)); err != nil {
				return err
			}
			input.Equipment[eq.Slot] = eq
		} else { // Assume misc options
			if input.Options == nil {
				input.Options = make(Options)
			}
			if err := input.Options.UnmarshalSimC([]byte(line)); err != nil {
				return err
			}
		}
	}
	return nil
}

// unmarshalField handles unmarshaling a single key=value into a field.
func unmarshalField(field reflect.Value, tag, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Ptr:
		if value == "" {
			return nil // Leave as nil
		}
		field.Set(reflect.New(field.Type().Elem()))
		return unmarshalField(field.Elem(), tag, value)
	case reflect.Slice:
		if value == "" {
			return nil
		}
		vals := strings.Split(value, "/")
		slice := reflect.MakeSlice(field.Type(), len(vals), len(vals))
		for i, v := range vals {
			elem := slice.Index(i)
			if err := unmarshalField(elem, "", v); err != nil {
				return err
			}
		}
		field.Set(slice)
	case reflect.Struct:
		if unmarshaler, ok := field.Addr().Interface().(Unmarshaler); ok {
			return unmarshaler.UnmarshalSimC([]byte(value))
		}
		return errors.New("nested struct without UnmarshalSimC")
	default:
		return fmt.Errorf("unsupported type for field %s", tag)
	}
	return nil
}
