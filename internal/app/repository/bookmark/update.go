package bookmark

import (
	"context"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
)

func (r *postgresRepository) Update(ctx context.Context, bookmark *model.Bookmark) error {
	result := r.db.WithContext(ctx).
		Model(&model.Bookmark{}).
		Where("id = ? AND user_id = ?", bookmark.ID, bookmark.UserID).
		Updates(map[string]any{
			"description": bookmark.Description,
			"url":         bookmark.URL,
		})
	if result.Error != nil {
		return dbutils.CatchErr(result.Error)
	}
	if result.RowsAffected == 0 {
		return dbutils.ErrRecordNotFound
	}

	return nil
}
