package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var errInvalidToken = errors.New("invalid token")

type jwtValidatorStub struct {
	validateTokenFn func(string) (jwt.MapClaims, error)
}

func (s *jwtValidatorStub) ValidateToken(token string) (jwt.MapClaims, error) {
	return s.validateTokenFn(token)
}

func TestJWTAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name             string
		authorization    string
		validator        *jwtValidatorStub
		expectedStatus   int
		expectedBody     string
		expectNextCalled bool
	}{
		{
			name:           "rejects missing authorization header",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Unauthorized"}`,
			validator: &jwtValidatorStub{
				validateTokenFn: func(string) (jwt.MapClaims, error) {
					t.Fatal("validator must not be called without an authorization header")
					return nil, nil
				},
			},
		},
		{
			name:          "rejects token validation error",
			authorization: "Bearer invalid-token",
			validator: &jwtValidatorStub{
				validateTokenFn: func(token string) (jwt.MapClaims, error) {
					assert.Equal(t, "invalid-token", token)
					return nil, errInvalidToken
				},
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Unauthorized"}`,
		},
		{
			name:          "accepts valid bearer token and stores claims",
			authorization: "Bearer valid-token",
			validator: &jwtValidatorStub{
				validateTokenFn: func(token string) (jwt.MapClaims, error) {
					assert.Equal(t, "valid-token", token)
					return jwt.MapClaims{"sub": "user-id"}, nil
				},
			},
			expectedStatus:   http.StatusNoContent,
			expectNextCalled: true,
		},
		{
			name:          "passes non bearer header to validator unchanged",
			authorization: "Token valid-token",
			validator: &jwtValidatorStub{
				validateTokenFn: func(token string) (jwt.MapClaims, error) {
					assert.Equal(t, "Token valid-token", token)
					return nil, errInvalidToken
				},
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Unauthorized"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.Use(NewJWTAuth(tc.validator).JWTAuth())

			nextCalled := false
			router.GET("/", func(ctx *gin.Context) {
				nextCalled = true
				claims, exists := ctx.Get("claims")
				assert.True(t, exists)
				assert.Equal(t, jwt.MapClaims{"sub": "user-id"}, claims)
				ctx.Status(http.StatusNoContent)
			})

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Set("Authorization", tc.authorization)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Equal(t, tc.expectNextCalled, nextCalled)
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, recorder.Body.String())
			}
		})
	}
}
