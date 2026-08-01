package user

import (
	"context"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
)

type CreateUserParam struct {
	Username string
	Email    string
	Password string
}

func (s *userService) CreateUser(ctx context.Context, param *CreateUserParam) (*model.User, error) {
	// hash password
	hashedPass, err := s.hasher.Hash(param.Password)
	if err != nil {
		return nil, err
	}

	// create model user
	user := &model.User{
		Username: param.Username,
		Password: hashedPass,
		Email:    param.Email,
	}

	// call repo to create user
	newUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.userRepo.GetUserByUsername(ctx, username)
}