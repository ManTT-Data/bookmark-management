package sqldb

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitMockDB(t *testing.T) *gorm.DB {
	dsn := ":memory:"
	gormDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		t.Fatalf("Failed to open mock DB: %v", err)
	}

	return gormDB
}