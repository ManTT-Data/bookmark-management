package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/ManTT-Data/bookmark-management/internal/app/service/mocks"
)

func TestShortenURL_ShortenLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.ShortenUrl

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "normal case - success",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/link/shorten",
					bytes.NewBuffer([]byte(`{"url":"https://blah.com","exp":60}`)),
				)
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				serviceMock.On("CreateShortenLink", ctx, "https://blah.com", int64(60)).Return("1234567", nil)
				return serviceMock
			},

			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"Shorten URL generated successfully!","code":"1234567"}`,
		},
		{
			name: "fail case - service failed",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/link/shorten",
					bytes.NewBuffer([]byte(`{"url":"https://blah.com","exp":60}`)),
				)
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				serviceMock.On("CreateShortenLink", ctx, "https://blah.com", int64(60)).Return("", testErr)
				return serviceMock
			},

			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"message":"Processing error"}`,
		},
		{
			name: "fail case - bad input",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/link/shorten",
					bytes.NewBuffer([]byte(`{"url":"not a link","exp":60}`)),
				)
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				return serviceMock
			},

			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"Input error","details":["Url is invalid (url)"]}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupRequest(gc)
			mockService := tc.setupMockService(gc)
			handler := NewShortenURL(mockService)
			handler.ShortenLink(gc)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestShortenURL_Redirect(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.ShortenUrl

		expectedStatus int
	}{
		{
			name: "normal case - success",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/link/redirect/1234567",
					nil,
				)
				ctx.Params = gin.Params{{Key: "code", Value: "1234567"}}
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				serviceMock.On("GetLinkFromCode", ctx, "1234567").Return("https://google.com", nil)
				return serviceMock
			},

			expectedStatus: http.StatusFound,
		},
		{
			name: "err case - wrong code",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/link/redirect",
					nil,
				)
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				return serviceMock
			},

			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "err case - service err",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/link/redirect/1234567",
					nil,
				)
				ctx.Params = gin.Params{{Key: "code", Value: "1234567"}}
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				serviceMock.On("GetLinkFromCode", ctx, "1234567").Return("", testErr)
				return serviceMock
			},

			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupRequest(gc)
			mockService := tc.setupMockService(gc)
			handler := NewShortenURL(mockService)
			handler.Redirect(gc)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			mockService.AssertExpectations(t)
		})
	}
}
