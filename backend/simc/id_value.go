package simc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/RobertsMJ/simc-cloud/backend/simc/internal"
	"github.com/samber/lo"
)

var ErrInvalidIDValueFormat = errors.New("invalid idvalue format")

// IDValueList represents a list of ID:Value pairs
// i.e. omnium_talents=136819:1/136822:1
type IDValueList []IDValue
type IDValue struct {
	ID    int
	Value int
}

var (
	_ internal.ValueMarshaler   = IDValueList(nil)
	_ internal.ValueUnmarshaler = (*IDValueList)(nil)
	_ internal.ValueMarshaler   = IDValue{}
	_ internal.ValueUnmarshaler = (*IDValue)(nil)
)

func (i *IDValue) UnmarshalStatement(s internal.Statement) error {
	if s.Key != "omnium_talents" {
		return ErrInvalidIDValueFormat
	}
	return i.UnmarshalSimc(s.Value)
}

func (i IDValue) MarshalSimc() (string, error) {
	return fmt.Sprintf("%d:%d", i.ID, i.Value), nil
}

func (i *IDValue) UnmarshalSimc(value string) error {
	k, v, found := strings.Cut(value, ":")
	if !found {
		return ErrInvalidIDValueFormat
	}
	id, err := strconv.Atoi(k)
	if err != nil {
		return fmt.Errorf("invalid ID in IDValue: %s", value)
	}
	val, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("invalid Value in IDValue: %s", value)
	}

	i.ID = id
	i.Value = val
	return nil
}

func (i *IDValueList) UnmarshalStatement(s internal.Statement) error {
	if s.Key != "omnium_talents" {
		return ErrInvalidIDValueFormat
	}
	return i.UnmarshalSimc(s.Value)
}

func (i IDValueList) MarshalSimc() (string, error) {
	res, err := lo.MapErr(i, func(idv IDValue, _ int) (string, error) {
		return idv.MarshalSimc()
	})
	if err != nil {
		return "", err
	}
	return strings.Join(res, "/"), nil
}

func (i *IDValueList) UnmarshalSimc(value string) error {
	parts := strings.Split(value, "/")
	res, err := lo.MapErr(parts, func(part string, _ int) (IDValue, error) {
		var idv IDValue
		err := idv.UnmarshalSimc(part)
		if err != nil {
			return IDValue{}, err
		}
		return idv, nil
	})
	if err != nil {
		return err
	}
	*i = IDValueList(res)
	return nil
}
