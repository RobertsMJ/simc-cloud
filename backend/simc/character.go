package simc

import (
	"encoding"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrMissingCharacterClass = fmt.Errorf("missing character class")
	ErrMissingCharacterName  = fmt.Errorf("missing character name")
	ErrMissingCharacterLevel = fmt.Errorf("missing character level")
	ErrMissingCharacterSpec  = fmt.Errorf("missing character spec")
	ErrMissingCharacterRace  = fmt.Errorf("missing character race")
)

const (
	charFieldLevel  = "level"
	charFieldRole   = "role"
	charFieldSpec   = "spec"
	charFieldRace   = "race"
	charFieldRegion = "region"
	charFieldServer = "server"
)

type Character struct {
	Name   string   `json:"name"`
	Class  WowClass `json:"class"`
	Level  int      `json:"level"`
	Role   string   `json:"role"`
	Spec   string   `json:"spec"`
	Race   string   `json:"race"`
	Region string   `json:"region"`
	Server string   `json:"server"`
}

func (c Character) render() string {
	return strings.Join([]string{
		fmt.Sprintf("%s=%s", c.Class, c.Name),
		fmt.Sprintf("level=%d", c.Level),
		fmt.Sprintf("role=%s", c.Role),
		fmt.Sprintf("spec=%s", c.Spec),
		fmt.Sprintf("race=%s", c.Race),
		fmt.Sprintf("region=%s", c.Region),
		fmt.Sprintf("server=%s", c.Server),
	}, "\n")
}

func isCharacterKey(key string) bool {
	if _, err := newWowClass(key); err == nil {
		return true
	}
	switch key {
	case charFieldLevel,
		charFieldRole,
		charFieldSpec,
		charFieldRace,
		charFieldRegion,
		charFieldServer:
		return true
	default:
		return false
	}
}

func newCharacter(statements []statement) (Character, error) {
	var c Character
	var err error

	for _, s := range statements {
		if wowclass, err := newWowClass(s.Key); err == nil {
			c.Class = wowclass
			c.Name = s.Value
			continue
		}

		switch s.Key {
		case charFieldLevel:
			c.Level, err = strconv.Atoi(s.Value)
		case charFieldRole:
			c.Role = s.Value
		case charFieldSpec:
			c.Spec = s.Value
		case charFieldRace:
			c.Race = s.Value
		case charFieldRegion:
			c.Region = s.Value
		case charFieldServer:
			c.Server = s.Value
		}
		if err != nil {
			return Character{}, err
		}
	}

	if c.Class == "" {
		return Character{}, ErrMissingCharacterClass
	}
	if c.Name == "" {
		return Character{}, ErrMissingCharacterName
	}
	if c.Level == 0 {
		return Character{}, ErrMissingCharacterLevel
	}
	if c.Spec == "" {
		return Character{}, ErrMissingCharacterSpec
	}
	if c.Race == "" {
		return Character{}, ErrMissingCharacterRace
	}

	return c, nil
}

type WowClass string

var _ encoding.TextMarshaler = (*WowClass)(nil)

func (c WowClass) MarshalText() ([]byte, error) {
	return []byte(c), nil
}

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

func newWowClass(s string) (WowClass, error) {
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
		return c, nil
	}
	return "", fmt.Errorf("invalid wow class: %s", s)
}
