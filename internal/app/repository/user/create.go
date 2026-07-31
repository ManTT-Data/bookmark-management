package user

import (
	"context"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
)

func (r *postgresRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *postgresRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	return &user, err
}