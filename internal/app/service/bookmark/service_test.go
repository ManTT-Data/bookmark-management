package bookmark

import (
	"context"
	"errors"
	"testing"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	bookmarkRepo "github.com/ManTT-Data/bookmark-management/internal/app/repository/bookmark"
	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
	"github.com/ManTT-Data/bookmark-management/pkg/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type repositoryStub struct {
	mock.Mock
}

func (s *repositoryStub) Create(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error) {
	args := s.Called(ctx, bookmark)
	result, _ := args.Get(0).(*model.Bookmark)
	return result, args.Error(1)
}

func (s *repositoryStub) ListByUserID(ctx context.Context, userID string, option bookmarkRepo.ListOption) ([]model.Bookmark, int64, error) {
	args := s.Called(ctx, userID, option)
	result, _ := args.Get(0).([]model.Bookmark)
	return result, args.Get(1).(int64), args.Error(2)
}

func (s *repositoryStub) GetByCode(ctx context.Context, code string) (*model.Bookmark, error) {
	args := s.Called(ctx, code)
	result, _ := args.Get(0).(*model.Bookmark)
	return result, args.Error(1)
}

func (s *repositoryStub) GetByIDAndUserID(ctx context.Context, id, userID string) (*model.Bookmark, error) {
	args := s.Called(ctx, id, userID)
	result, _ := args.Get(0).(*model.Bookmark)
	return result, args.Error(1)
}

func (s *repositoryStub) Update(ctx context.Context, bookmark *model.Bookmark) error {
	return s.Called(ctx, bookmark).Error(0)
}

func (s *repositoryStub) Delete(ctx context.Context, id, userID string) error {
	return s.Called(ctx, id, userID).Error(0)
}

type keyGeneratorStub struct {
	keys []string
}

func (s *keyGeneratorStub) GenerateKey(length int) string {
	key := s.keys[0]
	s.keys = s.keys[1:]
	return key
}

var _ utils.KeyGenerator = (*keyGeneratorStub)(nil)

func TestService_CreateRetriesExistingCode(t *testing.T) {
	repo := new(repositoryStub)
	keyGen := &keyGeneratorStub{keys: []string{"existing001", "available002"}}
	service := NewService(repo, keyGen)
	ctx := context.Background()

	repo.On("GetByCode", ctx, "existing001").Return(&model.Bookmark{}, nil).Once()
	repo.On("GetByCode", ctx, "available002").Return((*model.Bookmark)(nil), dbutils.ErrRecordNotFound).Once()
	repo.On("Create", ctx, mock.MatchedBy(func(bookmark *model.Bookmark) bool {
		return bookmark.Code == "available002" && bookmark.UserID == "user-1"
	})).Return(&model.Bookmark{Code: "available002"}, nil).Once()

	bookmark, err := service.Create(ctx, CreateParam{Description: "Google", URL: "https://google.com", UserID: "user-1"})
	require.NoError(t, err)
	require.Equal(t, "available002", bookmark.Code)
	repo.AssertExpectations(t)
}

func TestService_ListValidatesPaginationAndSort(t *testing.T) {
	service := NewService(new(repositoryStub), &keyGeneratorStub{})

	_, err := service.List(context.Background(), ListParam{UserID: "user-1", Page: -1})
	require.ErrorIs(t, err, ErrInvalidPagination)

	_, err = service.List(context.Background(), ListParam{UserID: "user-1", Sort: "created_at; DROP TABLE bookmarks"})
	require.ErrorIs(t, err, ErrInvalidSort)
}

func TestService_CreatePropagatesRepositoryError(t *testing.T) {
	repo := new(repositoryStub)
	keyGen := &keyGeneratorStub{keys: []string{"available002"}}
	service := NewService(repo, keyGen)
	ctx := context.Background()
	testErr := errors.New("database unavailable")

	repo.On("GetByCode", ctx, "available002").Return((*model.Bookmark)(nil), dbutils.ErrRecordNotFound).Once()
	repo.On("Create", ctx, mock.Anything).Return((*model.Bookmark)(nil), testErr).Once()

	_, err := service.Create(ctx, CreateParam{URL: "https://google.com", UserID: "user-1"})
	require.ErrorIs(t, err, testErr)
	repo.AssertExpectations(t)
}
