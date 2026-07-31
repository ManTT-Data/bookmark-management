package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

type hasher struct{}

func NewHasher() Hasher {
	return &hasher{}
}

var ErrPasswordHashFailed = errors.New("password hash failed")
var ErrPasswordNotMatch = errors.New("password not match")

func (h *hasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrPasswordHashFailed
	}
	return string(bytes), nil
}

func (h *hasher) Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}