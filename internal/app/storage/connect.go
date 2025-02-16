package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	mpg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	DSN          string `required:"true"`
	MaxOpenConns int    `default:"10"`
	MaxIdleConns int    `default:"10"`
}

func SetupDatabase(cfg Config) (*sqlx.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	return sqlx.NewDb(db, "postgres"), nil
}

func MigrateDB(db *sqlx.DB) error {
	driver, err := mpg.WithInstance(db.DB, &mpg.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://config/migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
