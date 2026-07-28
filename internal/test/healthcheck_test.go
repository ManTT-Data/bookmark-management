package test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/ManTT-Data/bookmark-management/internal/api"
	redisPkg "github.com/ManTT-Data/bookmark-management/pkg/redis"
)

var testErr = errors.New("test error")

func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder
		setupRedis    func(ctx context.Context) *redis.Client

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "normal case",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("GET", "/health-check", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},
			setupRedis: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},

			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"OK","instance_id":"1234","service_name":"test"}`,
		},
		{
			name: "err case",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("GET", "/health-check", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},
			setupRedis: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},

			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"Message":"Instance is not ready!"}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// Generate test config
			cfg := &api.Config{
				AppPort:     "8080",
				ServiceName: "test",
				InstanceID:  "1234",
			}

			// generate test redis
			mockRedis := tc.setupRedis(ctx)

			// Initialize the Gin router
			testRouter := api.New(gin.New(), cfg, mockRedis)
			recorder := tc.setupTestHTTP(testRouter)
			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.JSONEq(t, tc.expectedResponseBody, recorder.Body.String())
		})
	}
}
