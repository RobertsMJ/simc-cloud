package simc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

var ErrInvalidIDValueFormat = errors.New("invalid idvalue format")

// IDValue represents an ID:Value pair, such as 136819:1
type IDValue struct {
	ID    int `json:"id"`
	Value int `json:"value"`
}

func NewIDValue(s string) (IDValue, error) {
	k, v, found := strings.Cut(s, ":")
	if !found {
		return IDValue{}, ErrInvalidIDValueFormat
	}
	id, err := strconv.Atoi(k)
	if err != nil {
		return IDValue{}, fmt.Errorf("invalid ID in IDValue: %s", s)
	}
	val, err := strconv.Atoi(v)
	if err != nil {
		return IDValue{}, fmt.Errorf("invalid Value in IDValue: %s", s)
	}

	return IDValue{ID: id, Value: val}, nil
}

func (i IDValue) MarshalSimc() (string, error) {
	return fmt.Sprintf("%d:%d", i.ID, i.Value), nil
}

// IDValueList represents a list of ID:Value pairs
// i.e. omnium_talents=136819:1/136822:1
type IDValueList []IDValue

func NewIDValueList(value string) (IDValueList, error) {
	parts := strings.Split(value, "/")
	i := make([]IDValue, 0, len(parts))

	for _, part := range parts {
		idv, err := NewIDValue(part)
		if err != nil {
			return nil, err
		}
		i = append(i, idv)
	}

	return i, nil
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
