package simc

import (
	"fmt"
	"strings"
)

type Config struct {
	Precision float64 `json:"precision"`
}

func (c *Config) render() string {
	fields := []string{
		fmt.Sprintf("target_error=%.2f", c.Precision),
	}

	return strings.Join(fields, "\n")
}

const (
	SimPrecisionLow  float64 = 1.0
	SimPrecisionMed  float64 = 0.5
	SimPrecisionHigh float64 = 0.1
)
