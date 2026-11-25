package main

import (
	"errors"
	"os"

	"github.com/RobertsMJ/simc-cloud/backend/cmd/migrate/migrations"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/RobertsMJ/simc-cloud/backend/platform/postgresx"
	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/stdlib"
)

func run() error {
	log := platform.NewLogger()
	cfg, err := postgresx.NewConfig()
	if err != nil {
		log.Error("failed to load postgres config", "err", err)
		return err
	}

	pool, err := postgresx.NewPostgresConnection(cfg)
	if err != nil {
		log.Error("failed to connect to postgres", "err", err)
		return err
	}
	defer pool.Close()

	sourceDriver, err := iofs.New(migrations.Files, ".")
	if err != nil {
		log.Error("failed to create source driver", "err", err)
		return err
	}

	db := stdlib.OpenDBFromPool(pool)
	dbDriver, err := pgxmigrate.WithInstance(db, &pgxmigrate.Config{})
	if err != nil {
		log.Error("failed to create db driver", "err", err)
		return err
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "pgx5", dbDriver)
	if err != nil {
		log.Error("failed to create migrate instance", "err", err)
		return err
	}
	defer func() {
		if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
			log.Warn("failed to close migrate", "source_err", srcErr, "db_err", dbErr)
		}
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Error("migration failed", "err", err)
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
