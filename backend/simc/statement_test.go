package simc

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

var update = flag.Bool("update", false, "update expected test results")

type StatementTestSuite struct {
	suite.Suite
}

func TestStatementTestSuite(t *testing.T) {
	suite.Run(t, new(StatementTestSuite))
}

func (s *StatementTestSuite) TestParse() {

	input, err := os.ReadFile("fixtures/test_dh.simc")
	s.NoError(err)

	res, err := parse(strings.NewReader(string(input)))
	s.NoError(err)

	expected := s.goldenValue("test_dh_statements", res)
	s.ElementsMatch(res, expected)
}

func (s *StatementTestSuite) goldenValue(filename string, actual []statement) []statement {
	s.T().Helper()

	path := "fixtures/" + filename + ".json"
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		s.T().Fatal("could not open golden file")
	}
	defer f.Close()

	if *update {
		js, err := json.MarshalIndent(actual, "", "  ")
		if err != nil {
			s.T().Fatal("could not marshal golden dataset")
		}
		_, err = f.Write(js)
		if err != nil {
			s.T().Fatal("could not write golden file")
		}
		return actual
	}

	content, err := io.ReadAll(f)
	if err != nil {
		s.T().Fatal("could not read golden file")
	}
	var expected []statement
	err = json.Unmarshal(content, &expected)
	if err != nil {
		s.T().Fatal("could not unmarshal golden dataset")
	}
	return expected
}
