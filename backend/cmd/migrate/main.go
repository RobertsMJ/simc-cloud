package main

import (
	"context"
	"os"

	"github.com/RobertsMJ/simc-cloud/backend/cmd/migrate/migrations"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/RobertsMJ/simc-cloud/backend/platform/scyllax"
	"github.com/scylladb/gocqlx/v3/migrate"
)

func run() error {
	log := platform.NewLogger()
	cfg, err := scyllax.NewScyllaConfig()
	if err != nil {
		return err
	}

	session, err := scyllax.NewScyllaConnection(cfg)
	if err != nil {
		return err
	}

	if err := migrate.FromFS(context.Background(), *session, migrations.Files); err != nil {
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
