package simc

import (
	"fmt"
	"log/slog"
	"reflect"
	"sync"

	"github.com/RobertsMJ/simc-cloud/simc/internal"
)

type Character struct {
	Name   string   `json:"name" simc:"-"`
	Class  WowClass `json:"class" simc:"-"`
	Level  int      `json:"level" simc:"level"`
	Role   string   `json:"role" simc:"role"`
	Spec   string   `json:"spec" simc:"spec"`
	Race   string   `json:"race" simc:"race"`
	Region string   `json:"region" simc:"region"`
	Server string   `json:"server" simc:"server"`
}

func isCharacterKey(key string) bool {
	keys := getCharacterKeys()
	return keys[key]
}

var getCharacterKeys = sync.OnceValue(func() map[string]bool {
	keys := map[string]bool{
		// Class keys
		string(WowClassDH):      true,
		string(WowClassDK):      true,
		string(WowClassDruid):   true,
		string(WowClassEvoker):  true,
		string(WowClassHunter):  true,
		string(WowClassMage):    true,
		string(WowClassMonk):    true,
		string(WowClassPaladin): true,
		string(WowClassPriest):  true,
		string(WowClassRogue):   true,
		string(WowClassShaman):  true,
		string(WowClassWarlock): true,
		string(WowClassWarrior): true,
	}

	// Infer from struct tags
	typ := reflect.TypeFor[Character]()
	for f := range typ.Fields() {
		if tag := f.Tag.Get("simc"); tag != "" {
			simcTag := internal.NewSimcTag(tag)
			keys[simcTag.Key] = true
		}
	}

	return keys
})

var charFieldIdx = sync.OnceValue(internal.IndexSimcFields[Character])

func (c Character) MarshalSimcValue() (string, error) {
	panic("not impl")
}

func (c *Character) UnmarshalStatement(s internal.Statement) error {
	if c == nil {
		return fmt.Errorf("character unmarshal: %w", internal.ErrUnmarshalIntoNilPtr)
	}

	if !isCharacterKey(s.Key) {
		slog.Debug("attempted to assign unmapped statement to a character", "key", s.Key, "value", s.Value)
		return nil
	}

	if isWowClass(s.Key) {
		c.Class = WowClass(s.Key)
		c.Name = s.Value
		return nil
	}

	fields := charFieldIdx()
	idx, ok := fields[s.Key]
	if !ok {
		return fmt.Errorf("character unmarshal(%v): %w", s.Key, internal.ErrInvalidKeyForType)
	}
	field := reflect.ValueOf(c).Elem().Field(idx)

	if err := internal.AssignToField(field, s.Value); err != nil {
		return fmt.Errorf("character unmarshal into field(%v): %w", s.Key, err)
	}

	return nil
}

type WowClass string

const (
	WowClassDK      WowClass = "deathknight"
	WowClassDH      WowClass = "demonhunter"
	WowClassDruid   WowClass = "druid"
	WowClassEvoker  WowClass = "evoker"
	WowClassHunter  WowClass = "hunter"
	WowClassMage    WowClass = "mage"
	WowClassMonk    WowClass = "monk"
	WowClassPaladin WowClass = "paladin"
	WowClassPriest  WowClass = "priest"
	WowClassRogue   WowClass = "rogue"
	WowClassShaman  WowClass = "shaman"
	WowClassWarlock WowClass = "warlock"
	WowClassWarrior WowClass = "warrior"
)

func isWowClass(s string) bool {
	switch c := WowClass(s); c {
	case WowClassDK,
		WowClassDH,
		WowClassDruid,
		WowClassEvoker,
		WowClassHunter,
		WowClassMage,
		WowClassMonk,
		WowClassPaladin,
		WowClassPriest,
		WowClassRogue,
		WowClassShaman,
		WowClassWarlock,
		WowClassWarrior:
		return true
	}
	return false
}
