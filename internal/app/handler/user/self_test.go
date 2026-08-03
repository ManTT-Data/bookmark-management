package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	serviceMocks "github.com/ManTT-Data/bookmark-management/internal/app/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_GetSelfInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		claims         any
		setupMock      func(*serviceMocks.Service)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "returns self information for valid claims",
			claims: jwt.MapClaims{"sub": "user-id"},
			setupMock: func(serviceMock *serviceMocks.Service) {
				serviceMock.On("GetSelfInfo", mock.Anything, "user-id").Return(&model.User{
					ID:       "user-id",
					Username: "alice",
					Email:    "alice@example.com",
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Get self info successfully","data":{"id":"user-id","username":"alice","email":"alice@example.com"}}`,
		},
		{
			name:   "returns internal error when service fails",
			claims: jwt.MapClaims{"sub": "user-id"},
			setupMock: func(serviceMock *serviceMocks.Service) {
				serviceMock.On("GetSelfInfo", mock.Anything, "user-id").Return(nil, errServiceFailure).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Processing error"}`,
		},
		{
			name:           "rejects missing claims before calling service",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Invalid Token"}`,
		},
		{
			name:           "rejects malformed claims before calling service",
			claims:         "not-jwt-claims",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Invalid Token"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serviceMock := serviceMocks.NewService(t)
			if tc.setupMock != nil {
				tc.setupMock(serviceMock)
			}
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
			if tc.claims != nil {
				ctx.Set("claims", tc.claims)
			}
			handler := NewUserHandler(serviceMock)

			handler.GetSelfInfo(ctx)

			require.Equal(t, tc.expectedStatus, recorder.Code)
			require.JSONEq(t, tc.expectedBody, recorder.Body.String())
		})
	}
}
