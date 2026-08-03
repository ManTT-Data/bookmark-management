package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name             string
		method           string
		expectedStatus   int
		expectNextCalled bool
	}{
		{
			name:             "continues non preflight requests",
			method:           http.MethodGet,
			expectedStatus:   http.StatusNoContent,
			expectNextCalled: true,
		},
		{
			name:           "returns no content for options preflight",
			method:         http.MethodOptions,
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.Use(CORS())

			nextCalled := false
			router.Any("/", func(ctx *gin.Context) {
				nextCalled = true
				ctx.Status(http.StatusNoContent)
			})

			request := httptest.NewRequest(tc.method, "/", nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Equal(t, tc.expectNextCalled, nextCalled)
			assert.Equal(t, "*", recorder.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "true", recorder.Header().Get("Access-Control-Allow-Credentials"))
			assert.Contains(t, recorder.Header().Get("Access-Control-Allow-Headers"), "Authorization")
			assert.Contains(t, recorder.Header().Get("Access-Control-Allow-Methods"), http.MethodOptions)
		})
	}
}
