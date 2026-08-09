package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Bookmark represents a saved URL belonging to a user.
type Bookmark struct {
	ID          string    `json:"id" gorm:"primaryKey;column:id;type:varchar(50);not null"`
	Code        string    `json:"code" gorm:"column:code;type:varchar(255);unique;not null"`
	Description string    `json:"description" gorm:"column:description;type:varchar(255)"`
	URL         string    `json:"url" gorm:"column:url;type:varchar(255);not null"`
	UserID      string    `json:"-" gorm:"column:user_id;type:varchar(50);not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;not null"`
}

func (b *Bookmark) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.NewString()
	}

	return nil
}
