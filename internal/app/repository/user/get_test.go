package user

import (
	"testing"

	"github.com/ManTT-Data/bookmark-management/internal/test/data/fixtures"
	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepository_GetUserByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		inputID     string
		expectedErr error
		verify      func(t *testing.T, username string)
	}{
		{
			name:    "returns existing user",
			inputID: "9ba595b9-345e-4474-ad59-79a7442bc15e",
			verify: func(t *testing.T, username string) {
				assert.Equal(t, "Temaz", username)
			},
		},
		{
			name:        "returns record not found for unknown ID",
			inputID:     "missing-user-id",
			expectedErr: dbutils.ErrRecordNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := fixtures.NewFixture(t, &fixtures.UserCommonTestCase{})
			repository := NewRepository(db)

			user, err := repository.GetUserByID(t.Context(), tc.inputID)

			assert.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr != nil {
				assert.Nil(t, user)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, tc.inputID, user.ID)
			tc.verify(t, user.Username)
		})
	}
}
