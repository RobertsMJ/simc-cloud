package simc

import (
	"io"
)

// Some expected strings:
// waist=,id=244573,bonus_id=12214/13667/12497/12066/8960/12384/8791/13622/12667,content_tuning=3615,crafted_stats=49/32,crafting_quality=5
// server=illidan
// omnium_talents=136819:1/136822:1

type ValueMarshaler interface{ MarshalSimcValue() (string, error) }
type ValueUnmarshaler interface{ UnmarshalSimcValue(string) error }

type Parameter struct {
	key   string
	value string
}

type SimulationDocument struct {
	Character      Character
	EquippedItems  Loadout
	LoadoutOptions LoadoutOptions
	Config         SimConfig
	Extra          []Parameter
}

type documentBuilder struct {
}

func (db *documentBuilder) Build() *SimulationDocument {
	panic("not implemented")
}

func newDocumentBuilder(statements []statement) (*documentBuilder, error) {
	panic("not implemented")
}

func Unmarshal(doc io.Reader) (*SimulationDocument, error) {
	statements, err := parse(doc)
	if err != nil {
		return nil, err
	}

	builder, err := newDocumentBuilder(statements)
	if err != nil {
		return nil, err
	}

	return builder.Build(), nil
}

func Marshal(doc *SimulationDocument) (io.Reader, error) {
	panic("not implemented")
}
