package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"plexus-bff-service-go/internal/app/config"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func Open(ctx context.Context, cfg config.DatabaseConfig) (*Postgres, error) {
	if !cfg.Enabled() {
		return nil, nil
	}

	if cfg.Postgres.Host == "" || cfg.Postgres.Database == "" || cfg.Postgres.User == "" {
		return nil, fmt.Errorf("database integration is enabled but postgres host/database/user configuration is incomplete")
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("parse postgres connection string: %w", err)
	}

	poolConfig.MaxConns = cfg.Postgres.MaxConns
	poolConfig.MinConns = cfg.Postgres.MinConns
	poolConfig.MaxConnLifetime = time.Duration(cfg.Postgres.MaxConnLifetime) * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Postgres{pool: pool}, nil
}

func (p *Postgres) Pool() *pgxpool.Pool {
	if p == nil {
		return nil
	}
	return p.pool
}

func (p *Postgres) Close() {
	if p == nil || p.pool == nil {
		return
	}
	p.pool.Close()
}
