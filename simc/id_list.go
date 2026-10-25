package simc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/RobertsMJ/simc-cloud/simc/internal"
	"github.com/samber/lo"
)

type IDList []int // List of IDs joined by '/'

var (
	_ internal.ValueMarshaler   = IDList(nil)
	_ internal.ValueUnmarshaler = (*IDList)(nil)
)

func (i IDList) MarshalSimc() (string, error) {
	return strings.Join(lo.Map(i, func(id int, _ int) string {
		return strconv.Itoa(id)
	}), "/"), nil
}

func (i *IDList) UnmarshalSimc(value string) error {
	if i == nil {
		return fmt.Errorf("idlist unmarshal: %w", internal.ErrUnmarshalIntoNilPtr)
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
