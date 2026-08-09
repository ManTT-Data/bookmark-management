package bookmark

import (
	"context"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
)

func (r *postgresRepository) Delete(ctx context.Context, id, userID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Bookmark{})
	if result.Error != nil {
		return dbutils.CatchErr(result.Error)
	}
	if result.RowsAffected == 0 {
		return dbutils.ErrRecordNotFound
	}

	return nil
}
