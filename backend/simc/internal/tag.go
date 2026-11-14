package internal

import "strings"

type SimcTag struct {
	Key       string
	Omitempty bool
}

func NewSimcTag(tag string) SimcTag {
	res := SimcTag{}
	parts := strings.Split(tag, ",")
	if len(parts) == 0 {
		return res
	}

	res.Key = parts[0]
	if len(parts) > 1 {
		for _, part := range parts[1:] {
			if part == "omitempty" {
				res.Omitempty = true
			}
		}
	}
	return res
}
