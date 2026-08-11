package simc

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

type nameClass struct {
	key  WowClass
	name string
}

type Character struct {
	Name  string   `json:"name"`
	Class WowClass `json:"class"`
	Level int      `json:"level"`
	Role  string   `json:"role"`
	Spec  string   `json:"spec"`
}

var (
	_ ValueMarshaler   = Character{}
	_ ValueUnmarshaler = (*Character)(nil)
)

func (c Character) MarshalSimcValue() (string, error) {
	panic("not implemented")
}

func (c Character) UnmarshalSimcValue(value string) error {
	panic("not implemented")
}
