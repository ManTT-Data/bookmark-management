package sqldb

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestMigrateSQLDB(t *testing.T) {
	if os.Getenv("TEST_POSTGRES_MIGRATIONS") != "1" {
		t.Skip("set TEST_POSTGRES_MIGRATIONS=1 to run PostgreSQL migration integration tests")
	}

	cfg, err := newConfig("")
	if err != nil {
		t.Fatalf("load database config: %v", err)
	}

	originalDBName := cfg.DBName
	cfg.DBName = "postgres"
	adminDB, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		t.Fatalf("connect to PostgreSQL: %v", err)
	}

	testDBName := fmt.Sprintf("bookmark_migration_%s", strings.ReplaceAll(uuid.NewString(), "-", ""))
	if err := adminDB.Exec("CREATE DATABASE " + testDBName).Error; err != nil {
		t.Fatalf("create test database: %v", err)
	}

	t.Cleanup(func() {
		sqlDB, err := adminDB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}

		cfg.DBName = "postgres"
		cleanupDB, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
		if err != nil {
			t.Errorf("connect to PostgreSQL for cleanup: %v", err)
			return
		}
		defer func() {
			sqlDB, err := cleanupDB.DB()
			if err == nil {
				_ = sqlDB.Close()
			}
		}()

		if err := cleanupDB.Exec("DROP DATABASE IF EXISTS " + testDBName + " WITH (FORCE)").Error; err != nil {
			t.Errorf("drop test database: %v", err)
		}
	})

	cfg.DBName = testDBName
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		t.Fatalf("connect to test database: %v", err)
	}
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}()

	const migrationPath = "file://../../internal/migration"
	if err := MigrateSQLDB(db, migrationPath); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	if err := MigrateSQLDB(db, migrationPath); err != nil {
		t.Fatalf("rerun migrations: %v", err)
	}

	var version struct {
		Version uint
		Dirty   bool
	}
	if err := db.Raw("SELECT version, dirty FROM schema_migrations").Scan(&version).Error; err != nil {
		t.Fatalf("read migration version: %v", err)
	}
	if version.Version != 3 || version.Dirty {
		t.Fatalf("unexpected migration state: version=%d dirty=%t", version.Version, version.Dirty)
	}

	var columnCount int64
	if err := db.Raw(`
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_schema = 'public'
			AND table_name = 'users'
			AND column_name IN ('id', 'username', 'password', 'email')
	`).Scan(&columnCount).Error; err != nil {
		t.Fatalf("inspect users columns: %v", err)
	}
	if columnCount != 4 {
		t.Fatalf("users table has %d expected columns, want 4", columnCount)
	}

	var uniqueIndexCount int64
	if err := db.Raw(`
		SELECT COUNT(*)
		FROM pg_indexes
		WHERE schemaname = 'public'
			AND tablename = 'users'
			AND indexname IN ('uni_username', 'uni_email')
	`).Scan(&uniqueIndexCount).Error; err != nil {
		t.Fatalf("inspect users indexes: %v", err)
	}
	if uniqueIndexCount != 2 {
		t.Fatalf("users table has %d expected unique indexes, want 2", uniqueIndexCount)
	}

	var primaryKeyCount int64
	if err := db.Raw(`
		SELECT COUNT(*)
		FROM information_schema.table_constraints
		WHERE table_schema = 'public'
			AND table_name = 'users'
			AND constraint_type = 'PRIMARY KEY'
			AND constraint_name = 'users_pkey'
	`).Scan(&primaryKeyCount).Error; err != nil {
		t.Fatalf("inspect users primary key: %v", err)
	}
	if primaryKeyCount != 1 {
		t.Fatalf("users table has %d primary key constraints, want 1", primaryKeyCount)
	}

	var bookmarkColumnCount int64
	if err := db.Raw(`
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_schema = 'public'
			AND table_name = 'bookmarks'
			AND column_name IN ('id', 'code', 'description', 'url', 'user_id', 'created_at', 'updated_at')
	`).Scan(&bookmarkColumnCount).Error; err != nil {
		t.Fatalf("inspect bookmark columns: %v", err)
	}
	if bookmarkColumnCount != 7 {
		t.Fatalf("bookmarks table has %d expected columns, want 7", bookmarkColumnCount)
	}

	var bookmarkConstraintCount int64
	if err := db.Raw(`
		SELECT COUNT(*)
		FROM information_schema.table_constraints
		WHERE table_schema = 'public'
			AND table_name = 'bookmarks'
			AND constraint_name IN ('bookmark_pkey', 'uni_code', 'fk_user_id')
	`).Scan(&bookmarkConstraintCount).Error; err != nil {
		t.Fatalf("inspect bookmark constraints: %v", err)
	}
	if bookmarkConstraintCount != 3 {
		t.Fatalf("bookmarks table has %d expected constraints, want 3", bookmarkConstraintCount)
	}

	var bookmarkIndexCount int64
	if err := db.Raw(`
		SELECT COUNT(*)
		FROM pg_indexes
		WHERE schemaname = 'public'
			AND tablename = 'bookmarks'
			AND indexname = 'idx_bookmarks_user_created_at'
	`).Scan(&bookmarkIndexCount).Error; err != nil {
		t.Fatalf("inspect bookmark indexes: %v", err)
	}
	if bookmarkIndexCount != 1 {
		t.Fatalf("bookmarks table has %d expected list indexes, want 1", bookmarkIndexCount)
	}

	cfg.DBName = originalDBName
}
