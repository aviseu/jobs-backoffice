package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/aviseu/jobs/internal/app/storage"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

type config struct {
	DB  storage.Config
	Log struct {
		Level slog.Level `default:"info"`
	}
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})))

	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	// load environment variables
	slog.Info("loading environment variables...")
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}

	// configure logging
	slog.Info("configuring logging...")
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.Log.Level}))
	slog.SetDefault(log)

	// setup database
	slog.Info("setting up database...")
	db, err := storage.SetupDatabase(cfg.DB)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			slog.Error(fmt.Errorf("failed to close database connection: %w", err).Error())
		}
	}(db)

	// migrate db
	slog.Info("migrating database...")
	if err := storage.MigrateDB(db); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}
