package simc

import (
	"fmt"
	"strings"
)

type Config struct {
	Precision SimPrecision `json:"precision"`
}

func (c *Config) render() string {
	fields := []string{
		fmt.Sprintf("target_error=%s", c.Precision.render()),
	}

	return strings.Join(fields, "\n")
}

type SimPrecision int

const (
	SimPrecisionLow SimPrecision = iota
	SimPrecisionMed
	SimPrecisionHigh
)

func (s SimPrecision) render() string {
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
	return fmt.Sprintf("%.2f", targetError)
}
