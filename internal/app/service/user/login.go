package user

import (
	"context"
	"errors"
	"time"

	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
	"github.com/golang-jwt/jwt/v5"
)

const tokenDuration = time.Hour * 24

var ErrInvalidCredentials = errors.New("invalid credentials")

func (s *userService) Login(ctx context.Context, username string, password string) (string, error) {
	// check user exists with username or email
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	switch {
	case err == nil:
		break
	case errors.Is(err, dbutils.ErrRecordNotFound):
		return "", ErrInvalidCredentials
	default:
		return "", err
	}

	// compare password
	if !s.hasher.Compare(user.Password, password) {
		return "", ErrInvalidCredentials
	}

	// generate token
	tokenContent := jwt.MapClaims{
		"sub": user.ID,
		"exp":   time.Now().Add(tokenDuration).Unix(),
		"email": user.Email,
		"iat":   time.Now().Unix(),
	}

	tokenString, err := s.jwtGenerator.GenerateToken(tokenContent)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}