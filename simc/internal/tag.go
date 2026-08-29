package internal

import "strings"

type SimcTag struct {
	Key   string
	Value string
}

func NewSimcTag(tag string) SimcTag {
	parts := strings.Split(tag, ",")
	if len(parts) == 2 {
		return SimcTag{
			Key:   parts[0],
			Value: parts[1],
		}
	}
	return SimcTag{
		Key:   parts[0],
		Value: "",
	}
}
