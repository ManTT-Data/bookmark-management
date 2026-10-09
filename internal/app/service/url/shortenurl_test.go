package url

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	bookmarkRepo "github.com/ManTT-Data/bookmark-management/internal/app/repository/bookmark"
	"github.com/ManTT-Data/bookmark-management/internal/app/repository/mocks"
	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
	mocks2 "github.com/ManTT-Data/bookmark-management/pkg/utils/mocks"
)

var (
	redisTestErr = errors.New("test err")
	dbTestErr    = errors.New("db test err")
)

type mockBookmarkRepo struct {
	mock.Mock
}

func (m *mockBookmarkRepo) Create(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error) {
	args := m.Called(ctx, bookmark)
	result, _ := args.Get(0).(*model.Bookmark)
	return result, args.Error(1)
}

func (m *mockBookmarkRepo) ListByUserID(ctx context.Context, userID string, option bookmarkRepo.ListOption) ([]model.Bookmark, int64, error) {
	args := m.Called(ctx, userID, option)
	result, _ := args.Get(0).([]model.Bookmark)
	return result, args.Get(1).(int64), args.Error(2)
}

func (m *mockBookmarkRepo) GetByCode(ctx context.Context, code string) (*model.Bookmark, error) {
	args := m.Called(ctx, code)
	result, _ := args.Get(0).(*model.Bookmark)
	return result, args.Error(1)
}

func (m *mockBookmarkRepo) GetByIDAndUserID(ctx context.Context, id, userID string) (*model.Bookmark, error) {
	args := m.Called(ctx, id, userID)
	result, _ := args.Get(0).(*model.Bookmark)
	return result, args.Error(1)
}

func (m *mockBookmarkRepo) Update(ctx context.Context, bookmark *model.Bookmark) error {
	return m.Called(ctx, bookmark).Error(0)
}

func (m *mockBookmarkRepo) Delete(ctx context.Context, id, userID string) error {
	return m.Called(ctx, id, userID).Error(0)
}

func TestService_GetLinkFromCode(t *testing.T) {
	testCases := []struct {
		name string
		code string

		setupRepo         func(ctx context.Context) *mocks.URLStorage
		setupBookmarkRepo func(ctx context.Context) *mockBookmarkRepo

		expectedResult string
		expectedErr    error
	}{
		{
			name: "redis case - success (len 7)",
			code: "1234567",
			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mockRepo := mocks.NewURLStorage(t)
				mockRepo.On("GetURL", ctx, "1234567").Return("https://redis-link.com", nil)
				return mockRepo
			},
			setupBookmarkRepo: func(ctx context.Context) *mockBookmarkRepo {
				return nil
			},
			expectedResult: "https://redis-link.com",
			expectedErr:    nil,
		},
		{
			name: "redis case - code does not exist (len 7)",
			code: "1234567",
			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mockRepo := mocks.NewURLStorage(t)
				mockRepo.On("GetURL", ctx, "1234567").Return("", redis.Nil)
				return mockRepo
			},
			setupBookmarkRepo: func(ctx context.Context) *mockBookmarkRepo {
				return nil
			},
			expectedResult: "",
			expectedErr:    ErrCodeDoesntExist,
		},
		{
			name: "redis case - repo err (len 7)",
			code: "1234567",
			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mockRepo := mocks.NewURLStorage(t)
				mockRepo.On("GetURL", ctx, "1234567").Return("", redisTestErr)
				return mockRepo
			},
			setupBookmarkRepo: func(ctx context.Context) *mockBookmarkRepo {
				return nil
			},
			expectedResult: "",
			expectedErr:    redisTestErr,
		},
		{
			name: "db bookmark case - success (len 8)",
			code: "12345678",
			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				return nil
			},
			setupBookmarkRepo: func(ctx context.Context) *mockBookmarkRepo {
				m := new(mockBookmarkRepo)
				m.On("GetByCode", ctx, "12345678").Return(&model.Bookmark{URL: "https://bookmark-db.com"}, nil)
				return m
			},
			expectedResult: "https://bookmark-db.com",
			expectedErr:    nil,
		},
		{
			name: "db bookmark case - not found (len 8)",
			code: "12345678",
			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				return nil
			},
			setupBookmarkRepo: func(ctx context.Context) *mockBookmarkRepo {
				m := new(mockBookmarkRepo)
				m.On("GetByCode", ctx, "12345678").Return((*model.Bookmark)(nil), dbutils.ErrRecordNotFound)
				return m
			},
			expectedResult: "",
			expectedErr:    ErrCodeDoesntExist,
		},
		{
			name: "db bookmark case - repo err (len 8)",
			code: "12345678",
			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				return nil
			},
			setupBookmarkRepo: func(ctx context.Context) *mockBookmarkRepo {
				m := new(mockBookmarkRepo)
				m.On("GetByCode", ctx, "12345678").Return((*model.Bookmark)(nil), dbTestErr)
				return m
			},
			expectedResult: "",
			expectedErr:    dbTestErr,
		},
		{
			name: "invalid length - too short (len 4)",
			code: "test",
			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				return nil
			},
			setupBookmarkRepo: func(ctx context.Context) *mockBookmarkRepo {
				return nil
			},
			expectedResult: "",
			expectedErr:    ErrCodeDoesntExist,
		},
		{
			name: "invalid length - too long (len 9)",
			code: "123456789",
			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				return nil
			},
			setupBookmarkRepo: func(ctx context.Context) *mockBookmarkRepo {
				return nil
			},
			expectedResult: "",
			expectedErr:    ErrCodeDoesntExist,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			mockRepo := tc.setupRepo(ctx)
			mockBmRepo := tc.setupBookmarkRepo(ctx)

			var bmRepo bookmarkRepo.Repository
			if mockBmRepo != nil {
				bmRepo = mockBmRepo
			}

			testService := NewShortenUrl(mockRepo, bmRepo, nil)
			result, err := testService.GetLinkFromCode(ctx, tc.code)
			assert.Equal(t, tc.expectedResult, result)
			assert.ErrorIs(t, err, tc.expectedErr)

			if mockRepo != nil {
				mockRepo.AssertExpectations(t)
			}
			if mockBmRepo != nil {
				mockBmRepo.AssertExpectations(t)
			}
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

			testService := NewShortenUrl(mockRepo, nil, keygenMock)

			result, err := testService.CreateShortenLink(ctx, "https://blah.com", 60)
			assert.Equal(t, result, tc.expectedResult)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}

}
