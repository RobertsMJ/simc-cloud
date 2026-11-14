package simc

import (
	"fmt"
	"strings"

	"github.com/RobertsMJ/simc-cloud/backend/simc/internal"
)

type Config struct {
	Precision SimPrecision `json:"precision" simc:"target_error"`
}

func (c *Config) MarshalSimc() (string, error) {
	fields, err := internal.MarshalFields(c)
	if err != nil {
		return "", err
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
