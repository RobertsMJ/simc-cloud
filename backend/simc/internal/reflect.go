package internal

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
)

var (
	ErrInvalidStringValue  = errors.New("value is not a string")
	ErrNotAssignableField  = errors.New("field is not assignable")
	ErrNotAddressableField = errors.New("field is not addressable")
	ErrIncompatibleField   = errors.New("field is not compatible with value")
	ErrNilMarshal          = errors.New("t is nil")
)

func AssignToField(field reflect.Value, raw string) error {
	if !field.CanSet() {
		return ErrNotAssignableField
	}
	if !field.CanAddr() {
		return ErrNotAddressableField
	}

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

	val := reflect.ValueOf(raw)
	if unmarshaler, ok := field.Addr().Interface().(ValueUnmarshaler); ok {
		return unmarshaler.UnmarshalSimc(raw)
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

func MarshalFields(t any) ([]string, error) {
	v := reflect.ValueOf(t)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, ErrNilMarshal
		}
		v = v.Elem()
	}
	typ := v.Type()

	res := []string{}

	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := v.Field(i)

		if !fieldType.IsExported() {
			continue
		}

		tag := fieldType.Tag.Get("simc")
		if tag == "" || tag == "-" {
			continue
		}
		simcTag := NewSimcTag(tag)
		key := simcTag.Key

		if simcTag.Omitempty && fieldValue.IsZero() {
			continue
		}

		if fieldType.Type.Kind() == reflect.Slice && simcTag.Omitempty && fieldValue.Len() == 0 {
			continue
		}

		if fieldValue.Kind() == reflect.Pointer {
			if fieldValue.IsNil() {
				continue
			}
			fieldValue = fieldValue.Elem()
		}

		var val string
		if marshaler, ok := fieldValue.Interface().(ValueMarshaler); ok {
			var err error
			val, err = marshaler.MarshalSimc()
			if err != nil {
				return nil, fmt.Errorf("%s: %w", key, err)
			}
		} else {
			val = fmt.Sprintf("%v", fieldValue.Interface())
		}

		res = append(res, fmt.Sprintf("%s=%s", key, val))
	}
	return res, nil
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
