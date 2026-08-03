package url

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/ManTT-Data/bookmark-management/internal/app/repository/mocks"
	mocks2 "github.com/ManTT-Data/bookmark-management/pkg/utils/mocks"
)

var redisTestErr = errors.New("test err")

func TestService_GetLinkFromCode(t *testing.T) {
	testCases := []struct {
		name string

		setupRepo func(ctx context.Context) *mocks.URLStorage

		expectedResult string
		expectedErr    error
	}{
		{
			name: "normal case",

			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "test").Return("https://blah.com", nil)

				return mock
			},

			expectedResult: "https://blah.com",
			expectedErr:    nil,
		},
		{
			name: "empty case",

			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "test").Return("", redis.Nil)

				return mock
			},

			expectedResult: "",
			expectedErr:    ErrCodeDoesntExist,
		},
		{
			name: "err case",

			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "test").Return("", redisTestErr)

				return mock
			},

			expectedResult: "",
			expectedErr:    redisTestErr,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			mockRepo := tc.setupRepo(ctx)

			testService := NewShortenUrl(mockRepo, nil)
			result, err := testService.GetLinkFromCode(ctx, "test")
			assert.Equal(t, result, tc.expectedResult)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

var testErr = errors.New("test error")

const testExpTime = 60 * time.Second

func TestService_CreateShortenLink(t *testing.T) {
	testCases := []struct {
		name string

		setupRepo   func(ctx context.Context) *mocks.URLStorage
		setupKeyGen func() *mocks2.KeyGenerator

		expectedResult string
		expectedErr    error
	}{
		{
			name: "normal case - new key",

			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "1234567").Return("", redis.Nil)
				mock.On("StoreURL", ctx, "1234567", "https://blah.com", testExpTime).Return(nil)

				return mock
			},
			setupKeyGen: func() *mocks2.KeyGenerator {
				mockKeyGen := mocks2.NewKeyGenerator(t)
				mockKeyGen.On("GenerateKey", linkKeyLength).Return("1234567")

				return mockKeyGen
			},

			expectedResult: "1234567",
			expectedErr:    nil,
		},
		{
			name: "normal case - random the same key",

			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "1234567").Return("https://example.com", nil)
				mock.On("GetURL", ctx, "2345678").Return("", redis.Nil)
				mock.On("StoreURL", ctx, "2345678", "https://blah.com", testExpTime).Return(nil)

				return mock
			},
			setupKeyGen: func() *mocks2.KeyGenerator {
				mockKeyGen := mocks2.NewKeyGenerator(t)
				mockKeyGen.On("GenerateKey", linkKeyLength).Return("1234567").Once()
				mockKeyGen.On("GenerateKey", linkKeyLength).Return("2345678").Once()

				return mockKeyGen
			},

			expectedResult: "2345678",
			expectedErr:    nil,
		},
		{
			name: "err case - can't put new key",

			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "1234567").Return("", redis.Nil)
				mock.On("StoreURL", ctx, "1234567", "https://blah.com", testExpTime).Return(testErr)

				return mock
			},
			setupKeyGen: func() *mocks2.KeyGenerator {
				mockKeyGen := mocks2.NewKeyGenerator(t)
				mockKeyGen.On("GenerateKey", linkKeyLength).Return("1234567")

				return mockKeyGen
			},

			expectedResult: "",
			expectedErr:    testErr,
		},
		{
			name: "err case - can't get key",

			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "1234567").Return("", testErr)

				return mock
			},
			setupKeyGen: func() *mocks2.KeyGenerator {
				mockKeyGen := mocks2.NewKeyGenerator(t)
				mockKeyGen.On("GenerateKey", linkKeyLength).Return("1234567")

				return mockKeyGen
			},

			expectedResult: "",
			expectedErr:    testErr,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			mockRepo := tc.setupRepo(ctx)
			keygenMock := tc.setupKeyGen()

			testService := NewShortenUrl(mockRepo, keygenMock)

			result, err := testService.CreateShortenLink(ctx, "https://blah.com", 60)
			assert.Equal(t, result, tc.expectedResult)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}

}
