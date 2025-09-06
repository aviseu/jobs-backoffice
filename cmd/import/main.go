package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	ohttp "net/http"
	"os"
	"os/signal"
	"syscall"

	cpubsub "cloud.google.com/go/pubsub"
	"github.com/aviseu/jobs-backoffice/internal/app/application/http"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/importing"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/pubsub"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	"github.com/caarlos0/env/v11"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/net/netutil"
)

type config struct {
	PubSub struct {
		ProjectID  string        `env:"PROJECT_ID" required:"true"`
		JobTopicID string        `env:"JOB_TOPIC_ID,required"`
		Client     pubsub.Config `envPrefix:"CLIENT"`
	} `envPrefix:"PUBSUB_"`
	DB      storage.Config   `envPrefix:"DB_"`
	Import  http.Config      `envPrefix:"IMPORT_"`
	Gateway importing.Config `envPrefix:"GATEWAY_"`
	Log     struct {
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
	pjs := pubsub.NewJobService(client.Topic(cfg.PubSub.JobTopicID), cfg.PubSub.Client)

	// services
	slog.Info("setting up services...")
	chr := postgres.NewChannelRepository(db)
	ir := postgres.NewImportRepository(db)
	jr := postgres.NewJobRepository(db)

	is := importing.NewService(chr, ir, jr, ohttp.DefaultClient, cfg.Gateway, pjs, log)

	// start server
	server := http.SetupServer(ctx, cfg.Import, http.ImportRootHandler(is, log))
	listener, err := net.Listen("tcp", cfg.Import.Addr)
	if err != nil {
		return fmt.Errorf("failed to create listener on %s: %w", cfg.Import.Addr, err)
	}

	if cfg.Import.MaxConnections > 0 {
		listener = netutil.LimitListener(listener, cfg.Import.MaxConnections)
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("starting server...")
		serverErrors <- server.Serve(listener)
	}()

	// shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case <-done:
		slog.Info("shutting down server...")

		ctx, cancel := context.WithTimeout(ctx, cfg.Import.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
	}

	return nil
}
