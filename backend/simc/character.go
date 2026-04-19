package simc

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type WowClass string

const (
	DemonHunter WowClass = "demonhunter"
	DeathKnight WowClass = "deathknight"
	Druid       WowClass = "druid"
	Hunter      WowClass = "hunter"
	Mage        WowClass = "mage"
	Monk        WowClass = "monk"
	Paladin     WowClass = "paladin"
	Priest      WowClass = "priest"
	Rogue       WowClass = "rogue"
	Shaman      WowClass = "shaman"
	Warlock     WowClass = "warlock"
	Warrior     WowClass = "warrior"
	Evoker      WowClass = "evoker"
)

func validClassKeywords() map[string]struct{} {
	return map[string]struct{}{
		string(DemonHunter): {},
		string(DeathKnight): {},
		string(Druid):       {},
		string(Hunter):      {},
		string(Mage):        {},
		string(Monk):        {},
		string(Paladin):     {},
		string(Priest):      {},
		string(Rogue):       {},
		string(Shaman):      {},
		string(Warlock):     {},
		string(Warrior):     {},
		string(Evoker):      {},
	}
}

// Check if a string is a line defining character by checking for the presence of class= and a valid class keyword
func IsCharacterValue(keyword string) bool {
	kw := strings.Split(keyword, "=")[0] // Take the part before the first '=' as the keyword
	_, exists := validClassKeywords()[kw]
	if exists {
		return true
	}

	// Use reflection to check if the keyword matches any
	// of the simc struct tags in SimcCharacter
	charType := reflect.TypeOf(Character{})
	for field := range charType.Fields() {
		if field.Tag.Get("simc") == kw {
			return true
		}
	}

	return false
}

type Character struct {
	CharClass   WowClass `json:"class" simc:"keyword"`
	CharName    string   `json:"name" simc:"-"`
	Role        string   `json:"role" simc:"role"`
	Level       int      `json:"level" simc:"level"`
	Race        string   `json:"race" simc:"race"`
	Professions string   `json:"professions" simc:"professions"`
	Spec        string   `json:"spec" simc:"spec"`
	Talents     string   `json:"talents" simc:"talents"`
}

var _ Marshaler = (*Character)(nil)   // Ensure SimcCharacter implements SimCUnmarshaler
var _ Unmarshaler = (*Character)(nil) // Ensure SimcCharacter implements SimCUnmarshaler

func (c Character) MarshalSimC() ([]byte, error) {
	var sb strings.Builder
	if c.CharClass != "" && c.CharName != "" {
		sb.WriteString(string(c.CharClass) + "=" + strconv.Quote(c.CharName) + "\n") // Quote char name to handle spaces/special chars
	}

	val := reflect.ValueOf(c)
	typ := reflect.TypeOf(c)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("simc")
		if tag == "" || tag == "keyword" || tag == "-" {
			continue // Skip fields without simc tags or special tags
		}
		marshaled, err := marshalField(field, tag)
		if err != nil {
			return nil, err
		}
		if marshaled != nil {
			sb.Write(marshaled)
			sb.WriteString("\n")
		}
	}
	// Trim the final string to remove any trailing newlines
	return []byte(strings.TrimSpace(sb.String())), nil
}

// Iterates over each line of the data and maps the key=value pairs to the corresponding struct fields based on the simc tags
func (c *Character) UnmarshalSimC(data []byte) error {

	lines := strings.Split(string(data), "\n")
	val := reflect.ValueOf(c).Elem()
	typ := val.Type()
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// limit to 2 parts: key and value, don't try to split on '=' in the value
		// use case is something like professions=herbalism=88/mining=82
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return errors.New("invalid character line: " + line)
		}
		key := parts[0]
		value := parts[1]

		// If it's not a character value, skip it
		if !IsCharacterValue(key) {
			continue // Skip lines that don't match character fields
		}

		// Check if it's a class keyword
		if _, exists := validClassKeywords()[key]; exists {
			c.CharClass = WowClass(key)
			c.CharName = strings.Replace(value, "\"", "", -1) // Char names may be quoted
			continue
		}

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := typ.Field(i)
			tag := fieldType.Tag.Get("simc")
			if tag == key {
				switch field.Kind() {
				case reflect.String:
					field.SetString(value)
				case reflect.Int, reflect.Int64:
					intVal, err := strconv.Atoi(value)
					if err != nil {
						return err
					}
					field.SetInt(int64(intVal))
				case reflect.Struct:
					if unmarshaler, ok := field.Addr().Interface().(Unmarshaler); ok {
						if err := unmarshaler.UnmarshalSimC([]byte(line)); err != nil {
							return err
						}
					} else {
						return errors.New("nested struct without UnmarshalSimC")
					}
				default:
					return errors.New("unsupported field type in character struct")
				}
				break
			}
		}
	}
	return nil
}
