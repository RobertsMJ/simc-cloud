package scyllax

import (
	"fmt"

	"github.com/gocql/gocql"

	"github.com/scylladb/gocqlx/v3"
)

func NewScyllaConnection(cfg *ScyllaConfig) (*gocqlx.Session, error) {
	cluster := gocql.NewCluster(cfg.Hosts...)
	cluster.Keyspace = cfg.Keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.Username,
		Password: cfg.Password,
	}
	if cfg.DisableDiscovery {
		cluster.DisableInitialHostLookup = true
		cluster.IgnorePeerAddr = true
	}

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &session, nil
}
