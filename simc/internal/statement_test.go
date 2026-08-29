package internal_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/RobertsMJ/simc-cloud/simc/internal"
	"github.com/RobertsMJ/simc-cloud/test-utils"
	"github.com/stretchr/testify/suite"
)

type StatementTestSuite struct {
	suite.Suite
}

func TestStatementTestSuite(t *testing.T) {
	suite.Run(t, new(StatementTestSuite))
}

func (s *StatementTestSuite) TestParse() {

	input, err := os.ReadFile("fixtures/test_dh.simc")
	s.NoError(err)

	res, err := internal.ParseStatements(strings.NewReader(string(input)))
	s.NoError(err)

	path, err := filepath.Abs("./fixtures/test_dh_statements.json")
	s.NoError(err)
	expected := test.GoldenValue(s.T(), path, res)
	s.ElementsMatch(res, expected)
}
