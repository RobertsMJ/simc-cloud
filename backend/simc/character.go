package simc

type WowClass string

const (
	WowClassDH WowClass = "demonhunter"
)

type Character struct {
	Name        string      `json:"name"`
	Class       WowClass    `json:"class"`
	Level       int         `json:"level"`
	Race        string      `json:"race"`
	Server      string      `json:"server"`
	Role        string      `json:"role"`
	Professions Professions `json:"professions"`
	Spec        string      `json:"spec"`
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
