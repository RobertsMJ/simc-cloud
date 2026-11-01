package config

import "github.com/caarlos0/env/v11"

// Parse reads T from the environment using `env` struct tags. Every variable name is prefixed with prefix.
func Parse[T any](prefix string) (*T, error) {
	var cfg T
	if err := env.ParseWithOptions(&cfg, env.Options{Prefix: prefix}); err != nil {
		return nil, err
	}
	return &cfg, nil
}
