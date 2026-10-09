package dbutils

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var errFilter = []func(err error) (bool, error){
	FilterDuplicateEmail,
	FilterDuplicateUserName,
	FilterRecordNotFound,
}

var (
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUserName = errors.New("username already exists")
	ErrRecordNotFound    = errors.New("record not found")
)

func CatchErr(err error) error {
	if err == nil {
		return nil
	}

	for _, filter := range errFilter {
		if ok, e := filter(err); ok {
			return e
		}
	}
	return err
}

func FilterDuplicateEmail(err error) (bool, error) {
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "uni_email") || strings.Contains(message, "users.email") {
		return true, ErrDuplicateEmail
	}
	return false, nil
}

func FilterDuplicateUserName(err error) (bool, error) {
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "uni_username") || strings.Contains(message, "users.username") {
		return true, ErrDuplicateUserName
	}
	return false, nil
}
func FilterRecordNotFound(err error) (bool, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, ErrRecordNotFound
	}
	return false, nil
}
