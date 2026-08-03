package user

import (
	"testing"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	repositoryMocks "github.com/ManTT-Data/bookmark-management/internal/app/repository/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_GetUserByUsername(t *testing.T) {
	expectedUser := &model.User{ID: "user-id", Username: "alice"}

	testCases := []struct {
		name          string
		result        *model.User
		repositoryErr error
	}{
		{
			name:   "returns user from repository",
			result: expectedUser,
		},
		{
			name:          "returns repository error",
			repositoryErr: errRepository,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repositoryMock := repositoryMocks.NewRepository(t)
			repositoryMock.On("GetUserByUsername", mock.Anything, "alice").Return(tc.result, tc.repositoryErr).Once()
			service := NewService(repositoryMock, nil, nil)

			user, err := service.GetUserByUsername(t.Context(), "alice")

			require.Equal(t, tc.result, user)
			require.ErrorIs(t, err, tc.repositoryErr)
		})
	}
}
