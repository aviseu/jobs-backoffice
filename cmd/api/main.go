package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/aviseu/jobs/internal/app/domain/channel"
	"github.com/aviseu/jobs/internal/app/http"
	"github.com/aviseu/jobs/internal/app/storage"
	"github.com/aviseu/jobs/internal/app/storage/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

type config struct {
	DB  storage.Config
	API http.Config
	Log struct {
		Level slog.Level `default:"info"`
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
	r := postgres.NewChannelRepository(db)
	s := channel.NewService(r)

	// start server
	server := http.SetupServer(ctx, cfg.API, http.APIRootHandler(s, cfg.API, log))
	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("starting server...")
		serverErrors <- server.ListenAndServe()
	}()

	// shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case <-done:
		slog.Info("shutting down server...")

		ctx, cancel := context.WithTimeout(ctx, cfg.API.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
	}

	return nil
}
