package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/RobertsMJ/simc-cloud/backend/api"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/suite"
)

type ApiHandlerTestSuite struct {
	suite.Suite

	Api humatest.TestAPI
}

func TestApiHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ApiHandlerTestSuite))
}

func (suite *ApiHandlerTestSuite) SetupTest() {
	_, suite.Api = humatest.New(suite.T())
	api.RegisterRoutes(suite.Api, &api.ProfileHandler{})
}

func (suite *ApiHandlerTestSuite) TestPostProfile_Stub() {

	req := api.ProfileRequestBody{
		Profile: "test-profile",
	}

	expected := api.ProfileResponseBody{
		Profile: map[string]any{
			"req": map[string]any{
				"profile": "test-profile",
			},
		},
	}

	resp := suite.Api.Post("/api/profile", req)

	var actual api.ProfileResponseBody
	err := json.Unmarshal(resp.Body.Bytes(), &actual)
	suite.NoError(err, "should unmarshal response body")

	suite.Equal(http.StatusOK, resp.Code)
	suite.Equal(expected, actual)
}
