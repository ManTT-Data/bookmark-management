package healthcheck

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/ManTT-Data/bookmark-management/internal/app/service/mocks"
)

var testErr = errors.New("test err")

func TestHealthCheckHandler_HandleFunc(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.Healthcheck

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "normal case - success",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/health-check",
					nil,
				)
			},
			setupMockService: func(ctx context.Context) *mocks.Healthcheck {
				serviceMock := mocks.NewHealthcheck(t)
				serviceMock.On("HealthCheck", ctx).Return(&model.Message{
					Message:     "success",
					ServiceName: "bookmark",
					InstanceID:  "12345678",
				}, nil)

				return serviceMock
			},

			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"success","service_name":"bookmark","instance_id":"12345678"}`,
		},
		{
			name: "err case - service err",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/health-check",
					nil,
				)
			},
			setupMockService: func(ctx context.Context) *mocks.Healthcheck {
				serviceMock := mocks.NewHealthcheck(t)
				serviceMock.On("HealthCheck", ctx).Return(nil, testErr)

				return serviceMock
			},

			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"Message":"Instance is not ready!"}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupRequest(gc)
			mockService := tc.setupMockService(gc)
			handler := &healthCheckHandler{
				service: mockService,
			}
			handler.HealthCheck(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
