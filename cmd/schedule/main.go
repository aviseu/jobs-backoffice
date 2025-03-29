package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	cpubsub "cloud.google.com/go/pubsub"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/importing"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/pubsub"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

type config struct {
	PubSub struct {
		ProjectID     string `split_words:"true" required:"true"`
		ImportTopicID string `split_words:"true" required:"true"`
		Client        pubsub.Config
	} `split_words:"false"`
	DB  storage.Config
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

	// pubsub
	client, err := cpubsub.NewClient(ctx, cfg.PubSub.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to build pubsub client for project %s: %w", cfg.PubSub.ProjectID, err)
	}

	ps := pubsub.NewService(client.Topic(cfg.PubSub.ImportTopicID), cfg.PubSub.Client)

	// services
	slog.Info("setting up services...")
	ir := postgres.NewImportRepository(db)
	is := importing.NewImportService(ir)

	chr := postgres.NewChannelRepository(db)

	importSingle := importing.NewScheduleImportAction(is, ps, log)
	importActive := importing.NewScheduleImportsAction(chr, importSingle)

	slog.Info("starting imports...")
	if err := importActive.Execute(ctx); err != nil {
		return fmt.Errorf("failed to import active channels: %w", err)
	}

	slog.Info("all done.")

	return nil
}
