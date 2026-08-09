package bookmark

import (
	"context"
	"errors"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
)

func (s *service) Create(ctx context.Context, param CreateParam) (*model.Bookmark, error) {
	for range maxCodeAttempts {
		code := s.keyGen.GenerateKey(bookmarkCodeLength)
		_, err := s.repo.GetByCode(ctx, code)
		if err == nil {
			continue
		}
		if !errors.Is(err, dbutils.ErrRecordNotFound) {
			return nil, err
		}

		return s.repo.Create(ctx, &model.Bookmark{
			Code:        code,
			Description: param.Description,
			URL:         param.URL,
			UserID:      param.UserID,
		})
	}

	return nil, ErrCodeGeneration
}
