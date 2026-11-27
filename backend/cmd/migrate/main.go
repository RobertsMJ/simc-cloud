package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

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

	connCfg := pool.Config().ConnConfig
	log = log.With("host", connCfg.Host, "port", connCfg.Port, "database", connCfg.Database, "user", connCfg.User)

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

	migration, err := migrate.NewWithInstance("iofs", sourceDriver, "pgx5", dbDriver)
	if err != nil {
		log.Error("failed to create migrate instance", "err", err)
		return err
	}
	migration.Log = migrateLogger{log}
	defer func() {
		if srcErr, dbErr := migration.Close(); srcErr != nil || dbErr != nil {
			log.Warn("failed to close migrate", "source_err", srcErr, "db_err", dbErr)
		}
	}()

	err = migration.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Info("schema up to date")
		return nil
	}
	if err != nil {
		log.Error("migration failed", "err", err)
		return err
	}
	return nil
}

// migrateLogger adapts slog to golang-migrate's Logger, which reports each applied migration.
type migrateLogger struct{ log *slog.Logger }

func (l migrateLogger) Printf(format string, v ...any) {
	l.log.Info(strings.TrimSpace(fmt.Sprintf(format, v...)))
}

func (l migrateLogger) Verbose() bool {
	return l.log.Enabled(context.Background(), slog.LevelDebug)
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
