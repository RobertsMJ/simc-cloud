package simc

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) TestMarshal() {
	tc := []struct {
		Name      string
		Precision float64
		expected  string
	}{
		{Name: "marshals high precision", Precision: SimPrecisionHigh, expected: "target_error=0.10"},
		{Name: "marshals medium precision", Precision: SimPrecisionMed, expected: "target_error=0.50"},
		{Name: "marshals low precision", Precision: SimPrecisionLow, expected: "target_error=1.00"},
	}
	for _, tt := range tc {
		s.Run(tt.Name, func() {
			c := Config{Precision: tt.Precision}
			b := c.render()
			s.Equal(tt.expected, string(b))
		})
	}
}
