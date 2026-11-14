package simc_test

import (
	"testing"

	"github.com/RobertsMJ/simc-cloud/backend/simc"
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
		Precision simc.SimPrecision
		expected  string
	}{
		{Name: "marshals high precision", Precision: simc.SimPrecisionHigh, expected: "target_error=0.10"},
		{Name: "marshals medium precision", Precision: simc.SimPrecisionMed, expected: "target_error=0.50"},
		{Name: "marshals low precision", Precision: simc.SimPrecisionLow, expected: "target_error=1.00"},
	}
	for _, tt := range tc {
		s.Run(tt.Name, func() {
			c := simc.Config{Precision: tt.Precision}
			b := c.MarshalSimc()
			s.Equal(tt.expected, string(b))
		})
	}
}
