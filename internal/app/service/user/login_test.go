package user

import (
	"testing"
	"time"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	repositoryMocks "github.com/ManTT-Data/bookmark-management/internal/app/repository/mocks"
	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
	jwtMocks "github.com/ManTT-Data/bookmark-management/pkg/jwtutils/mocks"
	utilsMocks "github.com/ManTT-Data/bookmark-management/pkg/utils/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_Login(t *testing.T) {
	validUser := &model.User{
		ID:       "user-id",
		Username: "alice",
		Password: "hashed-password",
		Email:    "alice@example.com",
	}

	testCases := []struct {
		name          string
		setupMocks    func(*repositoryMocks.Repository, *utilsMocks.Hasher, *jwtMocks.JWTGenerator)
		expectedToken string
		expectedError error
	}{
		{
			name: "returns token for valid credentials",
			setupMocks: func(repositoryMock *repositoryMocks.Repository, hasherMock *utilsMocks.Hasher, jwtGeneratorMock *jwtMocks.JWTGenerator) {
				repositoryMock.On("GetUserByUsername", mock.Anything, "alice").Return(validUser, nil).Once()
				hasherMock.On("Compare", validUser.Password, "plain-password").Return(true).Once()
				jwtGeneratorMock.On("GenerateToken", mock.MatchedBy(func(claims jwt.MapClaims) bool {
					subject, subjectOK := claims["sub"].(string)
					email, emailOK := claims["email"].(string)
					issuedAt, issuedAtOK := claims["iat"].(int64)
					expiresAt, expiresAtOK := claims["exp"].(int64)

					return subjectOK && subject == validUser.ID &&
						emailOK && email == validUser.Email &&
						issuedAtOK && expiresAtOK && expiresAt > issuedAt &&
						expiresAt-issuedAt == int64((24*time.Hour).Seconds())
				})).Return("signed-token", nil).Once()
			},
			expectedToken: "signed-token",
		},
		{
			name: "maps missing user to invalid credentials",
			setupMocks: func(repositoryMock *repositoryMocks.Repository, _ *utilsMocks.Hasher, _ *jwtMocks.JWTGenerator) {
				repositoryMock.On("GetUserByUsername", mock.Anything, "alice").Return(nil, dbutils.ErrRecordNotFound).Once()
			},
			expectedError: ErrInvalidCredentials,
		},
		{
			name: "returns unexpected repository error",
			setupMocks: func(repositoryMock *repositoryMocks.Repository, _ *utilsMocks.Hasher, _ *jwtMocks.JWTGenerator) {
				repositoryMock.On("GetUserByUsername", mock.Anything, "alice").Return(nil, errRepository).Once()
			},
			expectedError: errRepository,
		},
		{
			name: "rejects invalid password",
			setupMocks: func(repositoryMock *repositoryMocks.Repository, hasherMock *utilsMocks.Hasher, _ *jwtMocks.JWTGenerator) {
				repositoryMock.On("GetUserByUsername", mock.Anything, "alice").Return(validUser, nil).Once()
				hasherMock.On("Compare", validUser.Password, "plain-password").Return(false).Once()
			},
			expectedError: ErrInvalidCredentials,
		},
		{
			name: "returns token generator error",
			setupMocks: func(repositoryMock *repositoryMocks.Repository, hasherMock *utilsMocks.Hasher, jwtGeneratorMock *jwtMocks.JWTGenerator) {
				repositoryMock.On("GetUserByUsername", mock.Anything, "alice").Return(validUser, nil).Once()
				hasherMock.On("Compare", validUser.Password, "plain-password").Return(true).Once()
				jwtGeneratorMock.On("GenerateToken", mock.Anything).Return("", errTokenCreation).Once()
			},
			expectedError: errTokenCreation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repositoryMock := repositoryMocks.NewRepository(t)
			hasherMock := utilsMocks.NewHasher(t)
			jwtGeneratorMock := jwtMocks.NewJWTGenerator(t)
			tc.setupMocks(repositoryMock, hasherMock, jwtGeneratorMock)
			service := NewService(repositoryMock, hasherMock, jwtGeneratorMock)

			token, err := service.Login(t.Context(), "alice", "plain-password")

			require.Equal(t, tc.expectedToken, token)
			require.ErrorIs(t, err, tc.expectedError)
		})
	}
}
