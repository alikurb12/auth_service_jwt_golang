package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxConns          = 25
	minConns          = 5
	maxConnLifetime   = time.Hour
	maxConnIdleTime   = 30 * time.Hour
	healthCheckPeriod = time.Minute
	connectTimeout    = 5 * time.Second
)

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse postgres config error: %v", err)
	}

	cfg.MaxConns = maxConns
	cfg.MinConns = minConns
	cfg.MaxConnLifetime = maxConnLifetime
	cfg.MaxConnIdleTime = maxConnIdleTime
	cfg.HealthCheckPeriod = healthCheckPeriod

	connectCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connectCtx, cfg)
	if err != nil {
		return nil, fmt.Errorf("Connection to db error: %v", err)
	}

	if err := pool.Ping(connectCtx); err != nil {
		return nil, fmt.Errorf("Pinging postgres: %v", err)
	}

	return pool, nil
}
