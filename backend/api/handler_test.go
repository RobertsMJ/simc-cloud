package api_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/RobertsMJ/simc-cloud/backend/api"
	"github.com/RobertsMJ/simc-cloud/backend/profile"
	"github.com/RobertsMJ/simc-cloud/backend/simc"
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
	profileSvc := profile.NewService()
	profileHandler := api.NewProfileHandler(profileSvc, slog.Default())
	api.RegisterRoutes(suite.Api, profileHandler)
}

func (suite *ApiHandlerTestSuite) TestPostProfile_Stub() {

	inProfile, err := os.ReadFile("testdata/kestw.simc")
	suite.NoError(err)

	doc, err := simc.NewSimulationDocument(bytes.NewReader(inProfile))
	expected := api.ProfileResponseBody{
		Profile: doc,
	}

	resp := suite.Api.Post("/profile", bytes.NewReader(inProfile))

	var actual api.ProfileResponseBody
	err = json.Unmarshal(resp.Body.Bytes(), &actual)
	suite.NoError(err, "should unmarshal response body")

	suite.Equal(http.StatusOK, resp.Code)
	suite.Equal(expected, actual)
}
