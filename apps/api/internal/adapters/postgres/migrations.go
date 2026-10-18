package postgres

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/001_initial_schema.sql
var initialSchemaSQL string

// Migrate executes versioned, idempotent PostgreSQL migrations.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin migration transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Ensure schema_migrations table exists
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`
	if _, err := tx.Exec(ctx, createTableSQL); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Check if version 1 is applied
	var count int
	err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = $1", 1).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to query migration status: %w", err)
	}

	if count == 0 {
		if _, err := tx.Exec(ctx, initialSchemaSQL); err != nil {
			return fmt.Errorf("failed to apply 001_initial_schema.sql: %w", err)
		}

		recordSQL := `INSERT INTO schema_migrations (version, name) VALUES ($1, $2)`
		if _, err := tx.Exec(ctx, recordSQL, 1, "001_initial_schema"); err != nil {
			return fmt.Errorf("failed to record migration 1: %w", err)
		}
	}

	return tx.Commit(ctx)
}
