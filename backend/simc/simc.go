package simc

import (
	"io"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/simc/internal"
)

func NewProfile(doc io.Reader) (*Profile, error) {
	p := &Profile{}

	statements, err := internal.ParseStatements(doc)
	if err != nil {
		return nil, err
	}

	for _, statement := range statements {
		switch {
		case isCharacterKey(statement.Key):
			if err := p.Character.UnmarshalStatement(statement); err != nil {
				slog.Error("could not unmarshal character statement", "err", err, "statement", statement)
				return nil, err
			}
		case statement.Key == "talents":
			if err = p.unmarshalTalent(statement); err != nil {
				slog.Error("could not handle talent statement", "err", err, "statement", statement)
				return nil, err
			}
		case isItemKey(statement.Key):
			if err = p.unmarshalItem(statement); err != nil {
				slog.Error("could not handle item statement", "err", err, "statement", statement)
				return nil, err
			}
		case statement.Key == "omnium_talents":
			if err = p.unmarshalOmniumTalents(statement); err != nil {
				slog.Error("could not handle omnium talents", "err", err, "statement", statement)
				return nil, err
			}
		default:
			if p.Extra == nil {
				p.Extra = []Parameter{}
			}
			p.Extra = append(p.Extra, Parameter{
				Key:   statement.Key,
				Value: statement.Value,
			})
		}
	}

	return p, nil
}
