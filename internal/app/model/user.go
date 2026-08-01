package model

import (
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type User struct {
	ID          string `json:"id" gorm:"primaryKey;column:id;type:varchar(50);not null"`
	Username    string `json:"username" gorm:"column:username;type:varchar(100);not null"`
	Password    string `json:"-" gorm:"column:password;type:varchar(255);not null"`
	Email       string `json:"email" gorm:"column:email;type:varchar(100);unique;not null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}