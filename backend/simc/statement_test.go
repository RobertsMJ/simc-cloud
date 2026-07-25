package simc_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type StatementTestSuite struct {
	suite.Suite
}

func TestStatementTestSuite(t *testing.T) {
	suite.Run(t, new(StatementTestSuite))
}
