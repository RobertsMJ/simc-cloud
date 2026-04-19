package apigw_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/RobertsMJ/simc-cloud-backend/simc"
	"github.com/RobertsMJ/simc-cloud-backend/transport/apigw"
	"github.com/stretchr/testify/suite"
)

type ApiGWTestSuite struct {
	suite.Suite
}

func TestApiGWTestSuite(t *testing.T) {
	suite.Run(t, new(ApiGWTestSuite))
}

func (s *ApiGWTestSuite) TestNewRequestHandler() {

	input, err := os.ReadFile("./test-data/event.json")
	s.NoError(err)

	var ev apigw.Request
	err = json.Unmarshal(input, &ev)
	s.NoError(err)

	handler := apigw.NewRequestHandler(func(ctx context.Context, req *simc.Input) (*simc.Output, error) {
		out := simc.Output("Hello, " + req.Character.CharName)
		return &out, nil
	})
	s.NotNil(handler)

	resp, err := handler(context.Background(), ev)
	s.NoError(err)
	s.Equal(200, resp.StatusCode)
	s.Equal("Hello, Kestdh", resp.Body)
}
