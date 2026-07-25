package simc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

var ErrInvalidProfessionFormat = errors.New("invalid profession format")

// Simc format for professions is wonky
type Professions []Profession // herbalism=88/mining=84
// Profession represents a single profession represented as ID=Value
// i.e. herbalism=88
type Profession struct {
	ID    string
	Value int
}

var (
	_ ValueMarshaler   = Professions(nil)
	_ ValueUnmarshaler = (*Professions)(nil)
	_ ValueMarshaler   = Profession{}
	_ ValueUnmarshaler = (*Profession)(nil)
)

func (p Professions) MarshalSimcValue() (string, error) {
	res, err := lo.MapErr(p, func(prof Profession, _ int) (string, error) {
		return prof.MarshalSimcValue()
	})
	if err != nil {
		return "", err
	}
	return strings.Join(res, "/"), nil
}

func (p *Professions) UnmarshalSimcValue(value string) error {
	parts := strings.Split(value, "/")
	res, err := lo.MapErr(parts, func(part string, _ int) (Profession, error) {
		var prof Profession
		err := prof.UnmarshalSimcValue(part)
		if err != nil {
			return Profession{}, err
		}
		return prof, nil
	})
	if err != nil {
		return err
	}
	*p = Professions(res)
	return nil
}

func (p Profession) MarshalSimcValue() (string, error) {
	return fmt.Sprintf("%s=%d", p.ID, p.Value), nil
}

func (p *Profession) UnmarshalSimcValue(value string) error {
	k, v, found := strings.Cut(value, "=")
	if !found {
		return ErrInvalidProfessionFormat
	}
	val, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("invalid Value in Profession: %s", value)
	}

	p.ID = k
	p.Value = val
	return nil
}
