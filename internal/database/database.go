package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/restaurantautomation/api/internal/config"
)

type Migration struct {
	Version   int64     `json:"version"`
	AppliedAt time.Time `json:"applied_at"`
}

type Inspector struct{ pool *pgxpool.Pool }

func NewInspector(pool *pgxpool.Pool) *Inspector { return &Inspector{pool: pool} }

func (inspector *Inspector) Migrations(ctx context.Context) ([]Migration, error) {
	rows, err := inspector.pool.Query(ctx, `SELECT version, applied_at FROM schema_migrations ORDER BY version DESC`)
	if err != nil {
		return nil, fmt.Errorf("list schema migrations: %w", err)
	}
	defer rows.Close()
	result := make([]Migration, 0)
	for rows.Next() {
		var migration Migration
		if err := rows.Scan(&migration.Version, &migration.AppliedAt); err != nil {
			return nil, fmt.Errorf("scan schema migration: %w", err)
		}
		result = append(result, migration)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate schema migrations: %w", err)
	}
	return result, nil
}

func Open(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("open database pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}
