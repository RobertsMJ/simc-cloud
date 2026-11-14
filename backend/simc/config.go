package simc

import "fmt"

type Config struct {
	Precision SimPrecision `json:"precision"`
}

type SimPrecision int

const (
	SimPrecisionLow SimPrecision = iota
	SimPrecisionMed
	SimPrecisionHigh
)

func (c *Config) MarshalSimc() string {
	var targetError float64
	switch c.Precision {
	case SimPrecisionHigh:
		targetError = 0.1
	case SimPrecisionMed:
		targetError = 0.5
	case SimPrecisionLow:
		fallthrough
	default:
		targetError = 1.0
	}
	return fmt.Sprintf("target_error=%.2f", targetError)
}
