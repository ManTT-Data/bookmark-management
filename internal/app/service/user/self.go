package user

import (
	"context"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
)

func (s *userService) GetSelfInfo(ctx context.Context, uid string) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}