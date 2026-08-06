package sqldb

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func MigrateSQLDB(db *gorm.DB, migrationPath string) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	pgDriver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	mig, err := migrate.NewWithDatabaseInstance(migrationPath, db.Name(), pgDriver)
	if err != nil {
		return err
	}

	err = mig.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return err
}

func migrateSchema(m *migrate.Migrate, mode string, steps int) error {
	var MigrationErr error

	switch mode {
	case "up":
		MigrationErr = m.Up()
	case "down":
		MigrationErr = m.Down()
	case "drop":
		MigrationErr = m.Drop()
	case "force":
		MigrationErr = m.Force(steps)
	case "step":
		if steps > 0 {
			MigrationErr = m.Steps(steps)
		} else {
			return errors.New("[DB Migration] steps must be positive integer")
		}
	default:
		return errors.New("[DB Migration]invalid mode")
	}

	if MigrationErr != nil && !errors.Is(MigrationErr, migrate.ErrNoChange) {
		return MigrationErr
	}

	return nil
}
