package simc

import (
	"strconv"
	"strings"
)

// IDList represents a list of IDs joined by '/'
type IDList []int

func (i IDList) render() string {
	ids := make([]string, 0, len(i))
	for _, id := range i {
		ids = append(ids, strconv.Itoa(id))
	}
	return strings.Join(ids, "/")
}

func newIDList(value string) (IDList, error) {
	i := make([]int, 0)

	if value == "" {
		return i, nil
	}

	for id := range strings.SplitSeq(value, "/") {
		parsedID, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		i = append(i, parsedID)
	}

	return i, nil
}
