package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/aviseu/jobs-backoffice/internal/app/domain"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/channel"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/job"
	"github.com/aviseu/jobs-backoffice/internal/app/gateway"
	"github.com/aviseu/jobs-backoffice/internal/app/storage"
	"github.com/aviseu/jobs-backoffice/internal/app/storage/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

type config struct {
	DB      storage.Config
	Gateway gateway.Config
	Log     struct {
		Level slog.Level `default:"info"`
	}
	Job struct {
		Workers int `default:"10"`
		Buffer  int `default:"10"`
	}
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})))

	if err := run(context.Background()); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
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

	// services
	slog.Info("setting up services...")
	ir := postgres.NewImportRepository(db)
	is := imports.NewService(ir)

	chr := postgres.NewChannelRepository(db)
	chs := channel.NewService(chr)

	jr := postgres.NewJobRepository(db)
	js := job.NewService(jr, cfg.Job.Buffer, cfg.Job.Workers)

	f := gateway.NewFactory(js, is, http.DefaultClient, cfg.Gateway, log)

	importActive := domain.NewImportActiveAction(chs, f, log)

	slog.Info("starting imports...")
	if err := importActive.Execute(ctx); err != nil {
		return fmt.Errorf("failed to import active channels: %w", err)
	}

	slog.Info("all done.")

	return nil
}
