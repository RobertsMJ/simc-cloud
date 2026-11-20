package simc

import (
	"strconv"
	"strings"

	"github.com/samber/lo"
)

// IDList represents a list of IDs joined by '/'
type IDList []int

func (i IDList) MarshalSimc() string {
	return strings.Join(lo.Map(i, func(id int, _ int) string {
		return strconv.Itoa(id)
	}), "/")
}

func NewIDList(value string) (IDList, error) {
	i := make([]int, 0)

	if value == "" {
		return i, nil
	}

	if res, err := lo.MapErr(
		strings.Split(value, "/"),
		func(s string, _ int) (int, error) {
			return strconv.Atoi(s)
		},
	); err != nil {
		return nil, err
	} else {
		i = append(i, res...)
	}
	return i, nil
}
