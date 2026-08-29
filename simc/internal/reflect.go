package internal

import (
	"errors"
	"log/slog"
	"reflect"
	"strconv"
)

var (
	ErrInvalidStringValue = errors.New("value is not a string")
	ErrNotAssignableField = errors.New("field is not assignable")
	ErrIncompatibleField  = errors.New("field is not compatible with value")
)

func AssignToField(field reflect.Value, raw string) error {
	if field.CanAddr() {
		if field.CanSet() {
			val := reflect.ValueOf(raw)
			if field.Kind() == reflect.Pointer ||
				field.Kind() == reflect.Slice ||
				field.Kind() == reflect.Struct {
				if field.IsNil() {
					switch field.Kind() {
					case reflect.Pointer:
						if raw != "" {
							field.Set(reflect.New(field.Type().Elem()))
							return AssignToField(field.Elem(), raw)
						}
					case reflect.Slice:
						field.Set(reflect.MakeSlice(reflect.SliceOf(field.Type().Elem()), 0, 1))
					case reflect.Struct:
						field.Set(reflect.New(field.Type().Elem()))
					}
				}
			}

			if unmarshaler, ok := field.Addr().Interface().(ValueUnmarshaler); ok {
				return unmarshaler.UnmarshalSimcValue(raw)
			} else if val.Type().ConvertibleTo(field.Type()) {
				field.Set(val.Convert(field.Type()))
				return nil
			} else {
				switch field.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					v, err := strconv.ParseInt(raw, 10, field.Type().Bits())
					if err != nil {
						slog.Debug("failed to parse int", "err", err, "raw", raw)
						return err
					}
					rValue := reflect.ValueOf(v)
					field.Set(rValue.Convert(field.Type()))
					return nil
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					v, err := strconv.ParseUint(raw, 10, field.Type().Bits())
					if err != nil {
						slog.Debug("failed to parse int", "err", err, "raw", raw)
						return err
					}
					rValue := reflect.ValueOf(v)
					field.Set(rValue.Convert(field.Type()))
					return nil
				case reflect.Float32, reflect.Float64:
					v, err := strconv.ParseFloat(raw, field.Type().Bits())
					if err != nil {
						slog.Debug("failed to parse float", "err", err, "raw", raw)
						return err
					}
					rValue := reflect.ValueOf(v)
					field.Set(rValue.Convert(field.Type()))
					return nil
				}
				return ErrIncompatibleField
			}
		}
	}
	return ErrNotAssignableField
}

func IndexSimcFields[T any]() map[string]int {
	fields := map[string]int{}
	typ := reflect.TypeFor[T]()
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if tag := f.Tag.Get("simc"); tag != "" && tag != "-" {
			simcTag := NewSimcTag(tag)
			fields[simcTag.Key] = i
		}
	}
	return fields
}
