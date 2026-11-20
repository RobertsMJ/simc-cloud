package simc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidIDValueFormat = errors.New("invalid idvalue format")

// IDValue represents an ID:Value pair, such as 136819:1
type IDValue struct {
	ID    int `json:"id"`
	Value int `json:"value"`
}

func newIDValue(s string) (IDValue, error) {
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

func (i IDValue) render() string {
	return fmt.Sprintf("%d:%d", i.ID, i.Value)
}

// IDValueList represents a list of ID:Value pairs
// i.e. omnium_talents=136819:1/136822:1
type IDValueList []IDValue

func newIDValueList(value string) (IDValueList, error) {
	parts := strings.Split(value, "/")
	i := make([]IDValue, 0, len(parts))

	for _, part := range parts {
		idv, err := newIDValue(part)
		if err != nil {
			return nil, err
		}
		i = append(i, idv)
	}

	return i, nil
}

func (i IDValueList) render() string {
	idv := make([]string, 0, len(i))
	for _, id := range i {
		idv = append(idv, id.render())
	}
	return strings.Join(idv, "/")
}
