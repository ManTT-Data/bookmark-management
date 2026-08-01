package fixtures

import (
	"testing"

	"github.com/ManTT-Data/bookmark-management/pkg/sqldb"
	"gorm.io/gorm"
)

type Fixture interface {
	SetupDB(db *gorm.DB)
	Migrate() error
	GenerateData() error
	GetDB() *gorm.DB
}

type BaseFixture struct {
	db *gorm.DB
}

func (f *BaseFixture) SetupDB(db *gorm.DB) {
	f.db = db
}

func (f *BaseFixture) DB() *gorm.DB {
	return f.db
}

func NewFixture(t *testing.T, fixture Fixture) *gorm.DB {
	//create test FB
	db := sqldb.InitMockDB(t)
	fixture.SetupDB(db)

	//mgrate schema
	err := fixture.Migrate()
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	//generate data
	err = fixture.GenerateData()
	if err != nil {
		t.Fatalf("Failed to generate data: %v", err)
	}

	//return db
	return fixture.GetDB()
}