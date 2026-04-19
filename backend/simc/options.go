package simc

import (
	"errors"
	"strings"
)

// TODO: At some point, type this more strongly
type Options map[string]string

var _ Marshaler = (*Options)(nil)   // Ensure SimcFightOptions implements SimCUnmarshaler
var _ Unmarshaler = (*Options)(nil) // Ensure SimcFightOptions implements SimCUnmarshaler

// Simply map key-value pairs
func (f Options) MarshalSimC() ([]byte, error) {
	var resString []string
	for key, value := range f {
		if key == "" {
			continue // Skip empty keys or values
		}
		resString = append(resString, key+"="+value)
	}
	return []byte(strings.Join(resString, "\n")), nil
}

func (f *Options) UnmarshalSimC(data []byte) error {
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			return errors.New("invalid fight option line: " + line)
		}
		key, value := kv[0], kv[1]
		(*f)[key] = value
	}
	return nil
}
