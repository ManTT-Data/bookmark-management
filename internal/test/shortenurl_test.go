package test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/ManTT-Data/bookmark-management/internal/api"
	redisPkg "github.com/ManTT-Data/bookmark-management/pkg/redis"
)

func TestShortenLinkEndpoint(t *testing.T) {
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
				req, _ := http.NewRequest("POST", "/v1/links/shorten", bytes.NewBuffer([]byte(`{"url":"https://blah.com","exp":60}`)))
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},
			setupRedis: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},

			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `"Shorten URL generated successfully!`,
		},
		{
			name: "err case",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("POST", "/v1/links/shorten", bytes.NewBuffer([]byte(`{"url":"https://blah.com","exp":1234}`)))
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
			expectedResponseBody: `{"message":"Processing error"}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// generate test redis
			mockRedis := tc.setupRedis(ctx)

			// Initialize the Gin router
			testRouter := api.New(gin.New(), &api.Config{}, mockRedis)
			recorder := tc.setupTestHTTP(testRouter)
			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)
		})
	}
}

func TestRedirectLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder
		setupRedis    func(ctx context.Context) *redis.Client

		expectedStatusCode int
	}{
		{
			name: "normal case",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("GET", "/v1/links/redirect/1234567", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},
			setupRedis: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				mock.Set(ctx, "1234567", "https://blah.com", 300)
				return mock
			},

			expectedStatusCode: http.StatusFound,
		},
		{
			name: "err case - code doesn't exist",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("GET", "/v1/links/redirect/1234567", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},
			setupRedis: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},

			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// generate test redis
			mockRedis := tc.setupRedis(ctx)

			// Initialize the Gin router
			testRouter := api.New(gin.New(), &api.Config{}, mockRedis)
			recorder := tc.setupTestHTTP(testRouter)
			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
		})
	}
}
