package test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	"github.com/ManTT-Data/bookmark-management/internal/api"
	"github.com/ManTT-Data/bookmark-management/internal/test/data/fixtures"
	jwtMocks "github.com/ManTT-Data/bookmark-management/pkg/jwtutils/mocks"
	redisPkg "github.com/ManTT-Data/bookmark-management/pkg/redis"
)

const fixtureUserID = "9ba595b9-345e-4474-ad59-79a7442bc15e"

var errInvalidFixtureToken = errors.New("invalid fixture token")

func TestSelfInfoEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		setupRequest   func(*http.Request)
		setupValidator func(*jwtMocks.JWTValidator)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "returns the fixture user for a valid token",
			setupRequest: func(request *http.Request) {
				request.Header.Set("Authorization", "Bearer fixture-token")
			},
			setupValidator: func(validatorMock *jwtMocks.JWTValidator) {
				validatorMock.On("ValidateToken", "fixture-token").Return(jwt.MapClaims{
					"sub": fixtureUserID,
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"message": "Get self info successfully",
				"data": {
					"id": "9ba595b9-345e-4474-ad59-79a7442bc15e",
					"username": "Temaz",
					"email": "test@example.com"
				}
			}`,
		},
		{
			name:           "rejects a missing authorization header",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Unauthorized"}`,
		},
		{
			name: "rejects an invalid token",
			setupRequest: func(request *http.Request) {
				request.Header.Set("Authorization", "Bearer invalid-token")
			},
			setupValidator: func(validatorMock *jwtMocks.JWTValidator) {
				validatorMock.On("ValidateToken", "invalid-token").Return(nil, errInvalidFixtureToken).Once()
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Unauthorized"}`,
		},
		{
			name: "returns processing error for an unknown authenticated user",
			setupRequest: func(request *http.Request) {
				request.Header.Set("Authorization", "Bearer unknown-user-token")
			},
			setupValidator: func(validatorMock *jwtMocks.JWTValidator) {
				validatorMock.On("ValidateToken", "unknown-user-token").Return(jwt.MapClaims{
					"sub": "unknown-user-id",
				}, nil).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Processing error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := fixtures.NewFixture(t, &fixtures.UserCommonTestCase{})
			redisMock := redisPkg.InitMockRedis(t)
			validatorMock := jwtMocks.NewJWTValidator(t)
			if tc.setupValidator != nil {
				tc.setupValidator(validatorMock)
			}

			engine := api.New(&api.EngineOpts{
				App:    gin.New(),
				Cfg:    &api.Config{},
				Redis:  redisMock,
				DB:     db,
				JWTVal: validatorMock,
			})
			request := httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
			if tc.setupRequest != nil {
				tc.setupRequest(request)
			}
			recorder := httptest.NewRecorder()

			engine.ServeHTTP(recorder, request)

			require.Equal(t, tc.expectedStatus, recorder.Code)
			require.JSONEq(t, tc.expectedBody, recorder.Body.String())
		})
	}
}
