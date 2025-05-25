package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v11"
	"log/slog"
	"os"

	cpubsub "cloud.google.com/go/pubsub"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/scheduling"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/pubsub"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type config struct {
	PubSub struct {
		ProjectID     string        `env:"PROJECT_ID" required:"true"`
		ImportTopicID string        `env:"IMPORT_TOPIC_ID,required"`
		Client        pubsub.Config `envPrefix:"CLIENT"`
	} `envPrefix:"PUBSUB_"`
	DB  storage.Config `envPrefix:"DB_"`
	Log struct {
		Level slog.Level `env:"LEVEL" envDefault:"info"`
	} `envPrefix:"LOG_"`
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
	if err := env.Parse(&cfg); err != nil {
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
	pis := pubsub.NewImportService(client.Topic(cfg.PubSub.ImportTopicID), cfg.PubSub.Client)

	// services
	slog.Info("setting up services...")
	chr := postgres.NewChannelRepository(db)
	ir := postgres.NewImportRepository(db)
	ss := scheduling.NewService(ir, chr, pis, log)

	slog.Info("starting imports...")
	if err := ss.ScheduleActiveChannels(ctx); err != nil {
		return fmt.Errorf("failed to import active channels: %w", err)
	}

	slog.Info("all done.")

	return nil
}
