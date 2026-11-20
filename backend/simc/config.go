package simc

import (
	"fmt"
	"strings"
)

type Config struct {
	Precision SimPrecision `json:"precision"`
}

func (c *Config) MarshalSimc() (string, error) {
	fields := []string{}

	if precisionStr, err := c.Precision.MarshalSimc(); err != nil {
		return "", err
	} else {
		fields = append(fields, fmt.Sprintf("target_error=%s", precisionStr))
	}

	return strings.Join(fields, "\n"), nil
}

type SimPrecision int

const (
	SimPrecisionLow SimPrecision = iota
	SimPrecisionMed
	SimPrecisionHigh
)

func (s SimPrecision) MarshalSimc() (string, error) {
	var targetError float64
	switch s {
	case SimPrecisionHigh:
		targetError = 0.1
	case SimPrecisionMed:
		targetError = 0.5
	case SimPrecisionLow:
		fallthrough
	default:
		targetError = 1.0
	}
	return fmt.Sprintf("%.2f", targetError), nil
}
