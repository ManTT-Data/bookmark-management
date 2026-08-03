package user

import (
	"testing"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	repositoryMocks "github.com/ManTT-Data/bookmark-management/internal/app/repository/mocks"
	utilsMocks "github.com/ManTT-Data/bookmark-management/pkg/utils/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_CreateUser(t *testing.T) {
	input := &CreateUserParam{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "plain-password",
	}

	testCases := []struct {
		name          string
		setupMocks    func(*repositoryMocks.Repository, *utilsMocks.Hasher)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name: "creates user with hashed password",
			setupMocks: func(repositoryMock *repositoryMocks.Repository, hasherMock *utilsMocks.Hasher) {
				hasherMock.On("Hash", input.Password).Return("hashed-password", nil).Once()
				repositoryMock.On("CreateUser", mock.Anything, &model.User{
					Username: input.Username,
					Email:    input.Email,
					Password: "hashed-password",
				}).Return(&model.User{
					Username: input.Username,
					Email:    input.Email,
					Password: "hashed-password",
				}, nil).Once()
			},
			expectedUser: &model.User{Username: input.Username, Email: input.Email, Password: "hashed-password"},
		},
		{
			name: "returns hashing error without creating user",
			setupMocks: func(_ *repositoryMocks.Repository, hasherMock *utilsMocks.Hasher) {
				hasherMock.On("Hash", input.Password).Return("", errHashing).Once()
			},
			expectedError: errHashing,
		},
		{
			name: "returns repository error",
			setupMocks: func(repositoryMock *repositoryMocks.Repository, hasherMock *utilsMocks.Hasher) {
				hasherMock.On("Hash", input.Password).Return("hashed-password", nil).Once()
				repositoryMock.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil, errRepository).Once()
			},
			expectedError: errRepository,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repositoryMock := repositoryMocks.NewRepository(t)
			hasherMock := utilsMocks.NewHasher(t)
			tc.setupMocks(repositoryMock, hasherMock)

			service := NewService(repositoryMock, hasherMock, nil)
			user, err := service.CreateUser(t.Context(), input)

			require.ErrorIs(t, err, tc.expectedError)
			require.Equal(t, tc.expectedUser, user)
		})
	}
}
