package bookmark

import (
	"context"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
)

func (r *postgresRepository) ListByUserID(ctx context.Context, userID string, option ListOption) ([]model.Bookmark, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Bookmark{}).Where("user_id = ?", userID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, dbutils.CatchErr(err)
	}

	var bookmarks []model.Bookmark
	if err := query.
		Order(bookmarkOrder(option.Sort)).
		Offset((option.Page - 1) * option.Limit).
		Limit(option.Limit).
		Find(&bookmarks).Error; err != nil {
		return nil, 0, dbutils.CatchErr(err)
	}

	return bookmarks, total, nil
}

func bookmarkOrder(sort string) string {
	switch sort {
	case "created_at":
		return "created_at DESC, id DESC"
	case "updated_at":
		return "updated_at DESC, id DESC"
	case "description":
		return "description ASC, id ASC"
	case "url":
		return "url ASC, id ASC"
	default:
		return "created_at DESC, id DESC"
	}
}
