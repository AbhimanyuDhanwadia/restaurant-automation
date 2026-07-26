package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type migration struct {
	version int64
	path    string
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, directory string) error {
	if _, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version BIGINT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`); err != nil {
		return fmt.Errorf("create migration ledger: %w", err)
	}
	migrations, err := readMigrations(directory)
	if err != nil {
		return err
	}
	for _, item := range migrations {
		var applied bool
		if err := pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`, item.version).Scan(&applied); err != nil {
			return fmt.Errorf("check migration %d: %w", item.version, err)
		}
		if applied {
			continue
		}
		sql, err := os.ReadFile(item.path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", item.path, err)
		}
		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin migration %d: %w", item.version, err)
		}
		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("apply migration %d: %w", item.version, err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, item.version); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("record migration %d: %w", item.version, err)
		}
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit migration %d: %w", item.version, err)
		}
	}
	return nil
}

func readMigrations(directory string) ([]migration, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("read migrations directory: %w", err)
	}
	result := make([]migration, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}
		prefix, _, ok := strings.Cut(entry.Name(), "_")
		if !ok {
			return nil, fmt.Errorf("migration %s must start with a numeric version", entry.Name())
		}
		version, err := strconv.ParseInt(prefix, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse migration %s: %w", entry.Name(), err)
		}
		result = append(result, migration{version: version, path: filepath.Join(directory, entry.Name())})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].version < result[j].version })
	for index := 1; index < len(result); index++ {
		if result[index].version == result[index-1].version {
			return nil, fmt.Errorf("duplicate migration version %d", result[index].version)
		}
	}
	return result, nil
}
