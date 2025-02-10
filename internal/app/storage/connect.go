package storage

import "github.com/jmoiron/sqlx"

type Config struct {
	DSN string `required:"true"`
}

func SetupDatabase(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	return db, nil
}
