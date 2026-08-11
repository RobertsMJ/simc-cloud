package simc

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/samber/lo"
)

// Some expected strings:
// waist=,id=244573,bonus_id=12214/13667/12497/12066/8960/12384/8791/13622/12667,content_tuning=3615,crafted_stats=49/32,crafting_quality=5
// server=illidan
// omnium_talents=136819:1/136822:1

type ValueMarshaler interface{ MarshalSimcValue() (string, error) }
type ValueUnmarshaler interface{ UnmarshalSimcValue(string) error }

type SimulationDocument struct {
	ActiveLoadout  Loadout
	LoadoutOptions LoadoutOptions
	Extra          []Parameter
}

type Parameter struct {
	Key   string
	Value string
}

func newParameter(s string) Parameter {
	k, v, found := strings.Cut(s, "=")
	if !found {
		return Parameter{Key: s, Value: ""}
	}
	return Parameter{Key: k, Value: v}
}

func NewSimulationDocument(doc io.Reader) (sd SimulationDocument, err error) {
	statements, err := parse(doc)
	if err != nil {
		return SimulationDocument{}, err
	}

	for _, statement := range statements {
		if statement.Key == "talent" {
			if err = sd.handleTalent(statement); err != nil {
				return sd, err
			}
		} else if _, ok := ItemSlotFromString(statement.Key); ok {
			if err = sd.handleItem(statement); err != nil {
				return sd, err
			}
		} else if statement.Key == "omnium_talents" {
			if err = sd.handleOmniumTalents(statement); err != nil {
				return sd, err
			}
		}
	}

	return sd, nil
}

func (sd *SimulationDocument) handleTalent(statement statement) error {
	if statement.Disabled {
		return sd.LoadoutOptions.UnmarshalStatement(statement)
	} else {
		return sd.ActiveLoadout.UnmarshalStatement(statement)
	}
}

func (sd *SimulationDocument) handleItem(statement statement) error {
	var item Item
	if err := UnmarshalStatement(statement, &item); err != nil {
		return err
	}

	if statement.Disabled {
		sd.LoadoutOptions.Items = append(sd.LoadoutOptions.Items, item)
	} else {
		if sd.ActiveLoadout.Items == nil {
			sd.ActiveLoadout.Items = make(map[ItemSlot]Item)
		}
		sd.ActiveLoadout.Items[item.Slot] = item
	}
	return nil
}

func (sd *SimulationDocument) handleOmniumTalents(statement statement) error {
	om := IDValueList{}
	if err := om.UnmarshalStatement(statement); err != nil {
		return err
	}

	if statement.Disabled {
		sd.LoadoutOptions.OmniumTalents = append(sd.LoadoutOptions.OmniumTalents, om)
	} else {
		sd.ActiveLoadout.OmniumTalents = om
	}
	return nil
}

func Marshal(doc *SimulationDocument) (io.Reader, error) {
	panic("not implemented")
}

// var index = sync.OnceValue[map[reflect.Type]reflect.]]

type simcTag struct {
	Key   string
	Value string
}

func newSimcTag(tag string) simcTag {
	parts := strings.Split(tag, ",")
	if len(parts) == 2 {
		return simcTag{
			Key:   parts[0],
			Value: parts[1],
		}
	}
	return simcTag{
		Key:   parts[0],
		Value: "",
	}
}

var reflectionIndex sync.Map // TODO:MJR Evaluate if this is the right approach

const (
	keyField   = "_key"
	valueField = "_value"
)

var (
	ErrNotPointer         = errors.New("value must be a pointer")
	ErrInvalidStringValue = errors.New("value is not a string")
	ErrNotAssignableField = errors.New("field is not assignable")
	ErrIncompatibleField  = errors.New("field is not compatible with value")
)

func assignToField(field reflect.Value, val any) error {
	if field.CanAddr() {
		if field.CanSet() {
			val := reflect.ValueOf(val)
			if unmarshaler, ok := field.Addr().Interface().(ValueUnmarshaler); ok {
				strVal, ok := val.Interface().(string)
				if !ok {
					return ErrInvalidStringValue
				}
				return unmarshaler.UnmarshalSimcValue(strVal)
			} else if val.Type().ConvertibleTo(field.Type()) {
				field.Set(val.Convert(field.Type()))
				return nil
			} else {
				return ErrIncompatibleField
			}
		}
	}
	return ErrNotAssignableField
}

func UnmarshalStatement(s statement, val any) error {
	t := reflect.TypeOf(val)
	if t.Kind() != reflect.Pointer {
		return ErrNotPointer
	}
	t = t.Elem()

	var indices map[string]int
	buildReflectionIndex(t)
	index, ok := reflectionIndex.Load(t)
	if !ok {
		panic("failed to load reflection index for type " + t.String())
	}
	if indices, ok = index.(map[string]int); !ok {
		panic("failed to cast reflection index for type " + t.String())
	}

	// Handle the key/value special case
	if keyIdx, ok := indices[keyField]; ok {
		field := reflect.ValueOf(val).Elem().Field(keyIdx)
		err := assignToField(field, s.Key)
		if err != nil {
			return err
		}
	}

	// Grab to values
	attrs := lo.Map(strings.Split(s.Value, ","), func(v string, _ int) string {
		return strings.TrimSpace(v)
	})

	// If we have a value before the first ',' it goes in the value field
	if len(attrs) > 0 && attrs[0] != "" {
		if valueIdx, ok := indices[valueField]; ok {
			attrVal := attrs[0]
			attrs = attrs[1:]

			field := reflect.ValueOf(val).Elem().Field(valueIdx)
			err := assignToField(field, attrVal)
			if err != nil {
				return err
			}
		}
	}

	// Filter out empty values
	vals := lo.Map(lo.Filter(attrs, func(v string, _ int) bool {
		return v != ""
	}), func(v string, _ int) Parameter {
		return newParameter(v)
	})
	// And map the parameters to the correspondingly tagged fields in the struct
	for _, v := range vals {
		if idx, ok := indices[v.Key]; ok {
			field := reflect.ValueOf(val).Elem().Field(idx)
			assignToField(field, v.Value)
		}
	}

	return nil
}

func buildReflectionIndex(t reflect.Type) map[string]int {
	if indices, ok := reflectionIndex.Load(t); ok {
		return indices.(map[string]int)
	}

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	fields := make(map[string]int)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("simc"); tag != "" {
			if tag == "-" {
				continue
			}

			simcTag := newSimcTag(tag)
			var key string
			switch simcTag.Value {
			case "key":
				key = keyField
			case "value":
				// TODO:MJR handle named key/value case
				key = valueField
			default:
				key = simcTag.Key
			}
			fields[key] = i
		}
	}

	reflectionIndex.Store(t, fields)
	return fields
}
