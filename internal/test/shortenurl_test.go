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
	"gorm.io/gorm"

	"github.com/ManTT-Data/bookmark-management/internal/api"
	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	redisPkg "github.com/ManTT-Data/bookmark-management/pkg/redis"
	"github.com/ManTT-Data/bookmark-management/pkg/sqldb"
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
			mockDB := sqldb.InitMockDB(t)

			// Initialize the Gin router
			testRouter := api.New(&api.EngineOpts{
				App:   gin.New(),
				Cfg:   &api.Config{},
				Redis: mockRedis,
				DB:    mockDB,
			})
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
		setupDB       func(t *testing.T, db *gorm.DB)

		expectedStatusCode int
		expectedLocation   string
	}{
		{
			name: "normal case - redis link (len 7)",
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
			expectedLocation:   "https://blah.com",
		},
		{
			name: "err case - redis code doesn't exist (len 7)",
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
		{
			name: "normal case - bookmark from db (len 8)",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("GET", "/v1/links/redirect/bm123456", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},
			setupRedis: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			setupDB: func(t *testing.T, db *gorm.DB) {
				err := db.AutoMigrate(&model.Bookmark{})
				assert.NoError(t, err)
				err = db.Create(&model.Bookmark{
					ID:     "bm-1",
					Code:   "bm123456",
					URL:    "https://google.com/search?q=golang",
					UserID: "user-1",
				}).Error
				assert.NoError(t, err)
			},
			expectedStatusCode: http.StatusFound,
			expectedLocation:   "https://google.com/search?q=golang",
		},
		{
			name: "err case - bookmark code doesn't exist in db (len 8)",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("GET", "/v1/links/redirect/nonexist", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},
			setupRedis: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			setupDB: func(t *testing.T, db *gorm.DB) {
				err := db.AutoMigrate(&model.Bookmark{})
				assert.NoError(t, err)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "err case - invalid code length (len 4)",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("GET", "/v1/links/redirect/test", nil)
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
			mockDB := sqldb.InitMockDB(t)
			if tc.setupDB != nil {
				tc.setupDB(t, mockDB)
			}

			// Initialize the Gin router
			testRouter := api.New(&api.EngineOpts{
				App:   gin.New(),
				Cfg:   &api.Config{},
				Redis: mockRedis,
				DB:    mockDB,
			})
			recorder := tc.setupTestHTTP(testRouter)
			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			if tc.expectedLocation != "" {
				assert.Equal(t, tc.expectedLocation, recorder.Header().Get("Location"))
			}
		})
	}
}
