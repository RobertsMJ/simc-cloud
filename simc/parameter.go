package simc

import "strings"

// Parameter contains an arbitrary key-value pair
type Parameter struct {
	Key   string
	Value string
}

func NewParameter(s string) Parameter {
	k, v, found := strings.Cut(s, "=")
	if !found {
		return Parameter{Key: s, Value: ""}
	}
	return Parameter{Key: k, Value: v}
}
