package user

import (
	"net/http"
	"testing"

	serviceMocks "github.com/ManTT-Data/bookmark-management/internal/app/service/mocks"
	service "github.com/ManTT-Data/bookmark-management/internal/app/service/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		body           string
		setupMock      func(*serviceMocks.Service)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "returns token for valid credentials",
			body: `{"username":"alice","password":"secret123"}`,
			setupMock: func(serviceMock *serviceMocks.Service) {
				serviceMock.On(
					"Login",
					mock.Anything,
					"alice",
					"secret123",
				).Return("signed-token", nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"signed-token"}`,
		},
		{
			name:           "rejects invalid input before calling service",
			body:           `{"username":"","password":""}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "maps invalid credentials to bad request",
			body: `{"username":"alice","password":"secret123"}`,
			setupMock: func(serviceMock *serviceMocks.Service) {
				serviceMock.On(
					"Login",
					mock.Anything,
					"alice",
					"secret123",
				).Return("", service.ErrInvalidCredentials).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Invalid credentials"}`,
		},
		{
			name: "maps unexpected service error to internal error",
			body: `{"username":"alice","password":"secret123"}`,
			setupMock: func(serviceMock *serviceMocks.Service) {
				serviceMock.On(
					"Login",
					mock.Anything,
					"alice",
					"secret123",
				).Return("", errServiceFailure).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Processing error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serviceMock := serviceMocks.NewService(t)
			if tc.setupMock != nil {
				tc.setupMock(serviceMock)
			}
			ctx, recorder := newJSONContext(tc.body)
			handler := NewUserHandler(serviceMock)

			handler.Login(ctx)

			require.Equal(t, tc.expectedStatus, recorder.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, recorder.Body.String())
			}
		})
	}
}
