package user

import "errors"

var (
	errHashing       = errors.New("hashing failed")
	errRepository    = errors.New("repository failed")
	errTokenCreation = errors.New("token creation failed")
)
