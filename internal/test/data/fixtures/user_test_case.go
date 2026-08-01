package fixtures

import (
	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"gorm.io/gorm"
)

type UserCommonTestCase struct {
	BaseFixture
}

func NewUserCommonTestCase() *UserCommonTestCase {
	return &UserCommonTestCase{}
}

func (u *UserCommonTestCase) SetupDB(db *gorm.DB) {
	u.BaseFixture.SetupDB(db)
}

func (u *UserCommonTestCase) Migrate() error {
	return u.db.AutoMigrate(&model.User{})
}

func (u *UserCommonTestCase) GenerateData() error {
	db := u.db.Session(&gorm.Session{SkipHooks: true})

	users := []model.User{
		{
			ID:       "9ba595b9-345e-4474-ad59-79a7442bc15e",
			Username: "Temaz",
			Password: "123456789",
			Email:    "test@example.com",
		},
		{
			ID:       "9ba595b2-345e-4474-ad59-79a7442bc25e",
			Username: "admin1",
			Password: "123456789",
			Email:    "test1@example.com",
		},
		{
			ID:       "9ba595b1-345e-4474-ad59-79a7442bc35e",
			Username: "admin2",
			Password: "123456789",
			Email:    "test2@example.com",
		},
	}

	return db.CreateInBatches(&users, 10).Error
}

func (u *UserCommonTestCase) GetDB() *gorm.DB {
	return u.db
}