package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

type IDList []int // List of IDs joined by '/'

var (
	_ ValueMarshaler   = IDList(nil)
	_ ValueUnmarshaler = (*IDList)(nil)
)

func (i *IDList) UnmarshalStatement(s Statement) error {
	panic("not implemented")
}

func (i IDList) MarshalSimcValue() (string, error) {
	return strings.Join(lo.Map(i, func(id int, _ int) string {
		return strconv.Itoa(id)
	}), "/"), nil
}

func (i *IDList) UnmarshalSimcValue(value string) error {
	if i == nil {
		return fmt.Errorf("idlist unmarshal: %w", ErrUnmarshalIntoNilPtr)
	}

	if value == "" {
		return nil
	}

	res, err := lo.MapErr(
		strings.Split(value, "/"),
		func(s string, _ int) (int, error) {
			return strconv.Atoi(s)
		},
	)
	if err != nil {
		return err
	}
	*i = append(*i, res...)
	return nil
}
