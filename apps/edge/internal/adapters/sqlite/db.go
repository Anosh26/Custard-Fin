package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// DB wraps the sql.DB instance with domain SQLite pragmas.
type DB struct {
	*sql.DB
}

// Open initializes and configures a SQLite connection pool with WAL mode and busy timeout.
func Open(ctx context.Context, dsn string) (*DB, error) {
	if dsn != ":memory:" && !isInMemory(dsn) {
		dir := filepath.Dir(dsn)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create sqlite database directory: %w", err)
			}
		}
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// SQLite single-writer configuration to eliminate database lock deadlocks
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// Apply robust performance and reliability PRAGMAs
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA busy_timeout = 5000;",
		"PRAGMA foreign_keys = ON;",
		"PRAGMA synchronous = NORMAL;",
	}

	for _, p := range pragmas {
		if _, err := db.ExecContext(ctx, p); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to set pragma %q: %w", p, err)
		}
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	wrapped := &DB{DB: db}
	if err := wrapped.Migrate(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run sqlite migrations: %w", err)
	}

	return wrapped, nil
}

func isInMemory(dsn string) bool {
	return dsn == "file::memory:?cache=shared" || dsn == ":memory:"
}
