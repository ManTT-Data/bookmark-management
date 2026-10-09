package user

import (
	"testing"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/ManTT-Data/bookmark-management/internal/test/data/fixtures"
	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSqlRepository_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB func(t *testing.T) *gorm.DB

		inputID       string
		inputUserName string
		inputPassword string
		inputEmail    string

		expectError bool
		verifyFunc  func(t *testing.T, db *gorm.DB, user *model.User, err error)
	}{
		{
			name: "Success",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestCase{})
			},
			inputUserName: "admin123",
			inputPassword: "123456789",
			inputEmail:    "user@gmail.com",
			expectError:   false,
			verifyFunc: func(t *testing.T, db *gorm.DB, user *model.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.NotEmpty(t, user.ID)

				savedUser := &model.User{}
				findErr := db.Where("username = ?", "admin123").First(savedUser).Error
				assert.NoError(t, findErr)
				assert.Equal(t, "admin123", savedUser.Username)
				assert.Equal(t, "123456789", savedUser.Password)
				assert.Equal(t, "user@gmail.com", savedUser.Email)
			},
		},
		{
			name: "Success_WithCustomID",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestCase{})
			},
			inputID:       "custom-uuid-001",
			inputUserName: "john_doe",
			inputPassword: "secretpassword",
			inputEmail:    "john@example.com",
			expectError:   false,
			verifyFunc: func(t *testing.T, db *gorm.DB, user *model.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "custom-uuid-001", user.ID)
			},
		},
		{
			name: "Error_DuplicateEmail",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixtures.NewFixture(t, &fixtures.UserCommonTestCase{})
				// Pre-insert a user with same email
				existingUser := &model.User{
					ID:       "existing-id-1",
					Username: "user1",
					Password: "password123",
					Email:    "duplicate@example.com",
				}
				db.Create(existingUser)
				return db
			},
			inputUserName: "user2",
			inputPassword: "password456",
			inputEmail:    "duplicate@example.com",
			expectError:   true,
			verifyFunc: func(t *testing.T, db *gorm.DB, user *model.User, err error) {
				assert.ErrorIs(t, err, dbutils.ErrDuplicateEmail)
				assert.Nil(t, user)
			},
		},
		{
			name: "Error_DuplicateID",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixtures.NewFixture(t, &fixtures.UserCommonTestCase{})

				existingUser := &model.User{
					ID:       "same-id-123",
					Username: "user_a",
					Password: "pass",
					Email:    "user_a@example.com",
				}
				db.Create(existingUser)
				return db
			},
			inputID:       "same-id-123",
			inputUserName: "user_b",
			inputPassword: "pass",
			inputEmail:    "user_b@example.com",
			expectError:   true,
			verifyFunc: func(t *testing.T, db *gorm.DB, user *model.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			r := NewRepository(db)
			db.AutoMigrate(&model.User{})

			user, err := r.CreateUser(ctx, &model.User{
				ID:       tc.inputID,
				Username: tc.inputUserName,
				Password: tc.inputPassword,
				Email:    tc.inputEmail,
			})

			if tc.verifyFunc != nil {
				tc.verifyFunc(t, db, user, err)
			}
		})
	}
}

func TestSqlRepository_GetUserByUsername(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB func(t *testing.T) *gorm.DB

		inputUsername string
		expectError   bool
		expectedErr   error
		verifyFunc    func(t *testing.T, user *model.User, err error)
	}{
		{
			name: "Success",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestCase{})
			},
			inputUsername: "Temaz",
			expectError:   false,
			verifyFunc: func(t *testing.T, user *model.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "9ba595b9-345e-4474-ad59-79a7442bc15e", user.ID)
				assert.Equal(t, "Temaz", user.Username)
				assert.Equal(t, "test@example.com", user.Email)
			},
		},
		{
			name: "Error_UserNotFound",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestCase{})
			},
			inputUsername: "non_existing_user",
			expectError:   true,
			expectedErr:   dbutils.ErrRecordNotFound,
			verifyFunc: func(t *testing.T, user *model.User, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, dbutils.ErrRecordNotFound)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			r := NewRepository(db)

			user, err := r.GetUserByUsername(ctx, tc.inputUsername)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			}
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, user, err)
			}
		})
	}
}
