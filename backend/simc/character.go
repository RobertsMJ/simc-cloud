package simc

type Character struct {
	Name        string
	Class       string
	Level       int
	Race        string
	Server      string
	Role        string
	Professions Professions
	Spec        string
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
