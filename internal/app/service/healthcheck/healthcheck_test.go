package healthcheck

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/ManTT-Data/bookmark-management/internal/app/repository/mocks"
)

func TestService_HealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		serviceName string

		expectedOutput *model.Message
		setupPingRepo  func(ctx context.Context) *mocks.Ping

		expectedErr error
	}{
		{
			name:        "Normal case 1",
			serviceName: "bookmark_app",

			setupPingRepo: func(ctx context.Context) *mocks.Ping {
				mock := mocks.NewPing(t)
				mock.On("Ping", ctx).Return(nil)

				return mock
			},

			expectedOutput: &model.Message{
				InstanceID:  "1234",
				Message:     "OK",
				ServiceName: "bookmark_app",
			},
			expectedErr: nil,
		},
		{
			name:        "Normal case 2",
			serviceName: "bookmark_service",

			setupPingRepo: func(ctx context.Context) *mocks.Ping {
				mock := mocks.NewPing(t)
				mock.On("Ping", ctx).Return(nil)

				return mock
			},

			expectedOutput: &model.Message{
				InstanceID:  "1234",
				Message:     "OK",
				ServiceName: "bookmark_service",
			},
			expectedErr: nil,
		},
		{
			name:        "Redis connection error",
			serviceName: "bookmark_app",

			setupPingRepo: func(ctx context.Context) *mocks.Ping {
				mock := mocks.NewPing(t)
				mock.On("Ping", ctx).Return(errRedisNotAvailable)

				return mock
			},

			expectedOutput: nil,
			expectedErr:    errRedisNotAvailable,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			redisMock := tc.setupPingRepo(ctx)

			testService := NewHealthcheck(tc.serviceName, "1234", redisMock)
			result, err := testService.HealthCheck(ctx)
			assert.Equal(t, result, tc.expectedOutput)
			assert.ErrorIs(t, err, tc.expectedErr)
			redisMock.AssertExpectations(t)
		})
	}

}
