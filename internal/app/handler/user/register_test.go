package user

import (
	"net/http"
	"testing"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	serviceMocks "github.com/ManTT-Data/bookmark-management/internal/app/service/mocks"
	service "github.com/ManTT-Data/bookmark-management/internal/app/service/user"
	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		body           string
		setupMock      func(*serviceMocks.Service)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "creates user from valid input",
			body: `{"username":"alice","email":"alice@example.com","password":"secret123"}`,
			setupMock: func(serviceMock *serviceMocks.Service) {
				serviceMock.On("CreateUser", mock.Anything, &service.CreateUserParam{
					Username: "alice",
					Email:    "alice@example.com",
					Password: "secret123",
				}).Return(&model.User{ID: "user-id", Username: "alice", Email: "alice@example.com"}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"user-id","username":"alice","email":"alice@example.com"}`,
		},
		{
			name:           "rejects invalid input before calling service",
			body:           `{"username":"","email":"invalid","password":"short"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "maps duplicate email to bad request",
			body: `{"username":"alice","email":"alice@example.com","password":"secret123"}`,
			setupMock: func(serviceMock *serviceMocks.Service) {
				serviceMock.On("CreateUser", mock.Anything, mock.AnythingOfType("*user.CreateUserParam")).Return(nil, dbutils.ErrDuplicateEmail).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Email already exists"}`,
		},
		{
			name: "maps duplicate username to bad request",
			body: `{"username":"alice","email":"alice@example.com","password":"secret123"}`,
			setupMock: func(serviceMock *serviceMocks.Service) {
				serviceMock.On("CreateUser", mock.Anything, mock.AnythingOfType("*user.CreateUserParam")).Return(nil, dbutils.ErrDuplicateUserName).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Username already exists"}`,
		},
		{
			name: "maps unexpected service error to internal error",
			body: `{"username":"alice","email":"alice@example.com","password":"secret123"}`,
			setupMock: func(serviceMock *serviceMocks.Service) {
				serviceMock.On("CreateUser", mock.Anything, mock.AnythingOfType("*user.CreateUserParam")).Return(nil, errServiceFailure).Once()
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

			handler.Register(ctx)

			require.Equal(t, tc.expectedStatus, recorder.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, recorder.Body.String())
			}
		})
	}
}
