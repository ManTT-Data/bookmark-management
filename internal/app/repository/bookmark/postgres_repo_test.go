package bookmark

import (
	"context"
	"testing"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
	"github.com/ManTT-Data/bookmark-management/pkg/sqldb"
	"github.com/stretchr/testify/require"
)

const testUserID = "user-1"

func TestPostgresRepository_ListByUserID(t *testing.T) {
	db := sqldb.InitMockDB(t)
	require.NoError(t, db.AutoMigrate(&model.Bookmark{}))

	repo := NewRepository(db)
	ctx := context.Background()
	bookmarks := []*model.Bookmark{
		{ID: "bookmark-1", Code: "bookmarkA1", URL: "https://a.example", Description: "A", UserID: testUserID},
		{ID: "bookmark-2", Code: "bookmarkB2", URL: "https://b.example", Description: "B", UserID: testUserID},
		{ID: "bookmark-3", Code: "bookmarkC3", URL: "https://c.example", Description: "C", UserID: "user-2"},
	}
	for _, bookmark := range bookmarks {
		_, err := repo.Create(ctx, bookmark)
		require.NoError(t, err)
	}

	result, total, err := repo.ListByUserID(ctx, testUserID, ListOption{Page: 1, Limit: 1, Sort: "created_at"})
	require.NoError(t, err)
	require.EqualValues(t, 2, total)
	require.Len(t, result, 1)
	require.Equal(t, testUserID, result[0].UserID)
}

func TestPostgresRepository_ScopesMutationToUser(t *testing.T) {
	db := sqldb.InitMockDB(t)
	require.NoError(t, db.AutoMigrate(&model.Bookmark{}))

	repo := NewRepository(db)
	ctx := context.Background()
	_, err := repo.Create(ctx, &model.Bookmark{ID: "bookmark-1", Code: "bookmarkA1", URL: "https://a.example", UserID: testUserID})
	require.NoError(t, err)

	err = repo.Update(ctx, &model.Bookmark{ID: "bookmark-1", UserID: "user-2", URL: "https://updated.example"})
	require.ErrorIs(t, err, dbutils.ErrRecordNotFound)

	err = repo.Delete(ctx, "bookmark-1", "user-2")
	require.ErrorIs(t, err, dbutils.ErrRecordNotFound)

	bookmark, err := repo.GetByIDAndUserID(ctx, "bookmark-1", testUserID)
	require.NoError(t, err)
	require.Equal(t, "https://a.example", bookmark.URL)
}
