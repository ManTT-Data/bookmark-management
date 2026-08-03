package url

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	redisPkg "github.com/ManTT-Data/bookmark-management/pkg/redis"
)

var testErr = errors.New("test error")

func TestUrlStorage_GetURL(t *testing.T) {
	testCases := []struct {
		name string

		setupMock func(ctx context.Context) *redis.Client

		expectedResult string
		expectedErr    error
	}{
		{
			name: "normal case",

			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				mock.Set(ctx, "test", "passed", 300)
				return mock
			},

			expectedResult: "passed",
			expectedErr:    nil,
		},
		{
			name: "err case - redis connection err",

			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},

			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			redisMock := tc.setupMock(ctx)

			testRepo := NewURLStorage(redisMock)
			result, err := testRepo.GetURL(ctx, "test")
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, result, tc.expectedResult)
		})
	}
}

func TestUrlStorage_StoreURL(t *testing.T) {
	testCases := []struct {
		name string

		setupMock func(ctx context.Context) *redis.Client

		expectedErr error
		verifyFunc  func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "normal case",

			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},

			expectedErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				assert.Nil(t, r.Get(ctx, "test").Err())
			},
		},
		{
			name: "err case - redis connection err",

			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},

			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			redisMock := tc.setupMock(ctx)

			testRepo := NewURLStorage(redisMock)
			err := testRepo.StoreURL(ctx, "test", "passed", 1)
			assert.ErrorIs(t, err, tc.expectedErr)
			if err == nil {
				tc.verifyFunc(ctx, redisMock)
			}
		})
	}
}
