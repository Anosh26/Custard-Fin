package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type migration struct {
	version int
	name    string
	sql     string
}

var migrations = []migration{
	{
		version: 1,
		name:    "create_media_table",
		sql: `
		CREATE TABLE IF NOT EXISTS media (
			id TEXT PRIMARY KEY,
			path TEXT UNIQUE NOT NULL,
			file_size INTEGER NOT NULL,
			mod_time DATETIME NOT NULL,
			content_hash TEXT NOT NULL,
			title TEXT NOT NULL,
			duration_s INTEGER NOT NULL,
			width INTEGER NOT NULL,
			height INTEGER NOT NULL,
			video_codec TEXT NOT NULL,
			audio_codec TEXT NOT NULL,
			is_public BOOLEAN NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_media_content_hash ON media(content_hash);
		CREATE INDEX IF NOT EXISTS idx_media_path ON media(path);
		`,
	},
	{
		version: 2,
		name:    "create_watch_progress_table",
		sql: `
		CREATE TABLE IF NOT EXISTS watch_progress (
			media_id TEXT PRIMARY KEY REFERENCES media(id) ON DELETE CASCADE,
			position_s INTEGER NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT 0,
			last_watched_at DATETIME NOT NULL
		);
		`,
	},
}

// Migrate applies all pending schema migrations inside a transaction.
func (db *DB) Migrate(ctx context.Context) error {
	// Ensure migration table exists
	initSQL := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at DATETIME NOT NULL
	);`
	if _, err := db.ExecContext(ctx, initSQL); err != nil {
		return fmt.Errorf("failed to initialize schema_migrations table: %w", err)
	}

	for _, m := range migrations {
		var exists int
		err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = ?", m.version).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed checking migration status for version %d: %w", m.version, err)
		}

		if exists > 0 {
			continue // already applied
		}

		tx, err := db.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", m.version, err)
		}

		if _, err := tx.ExecContext(ctx, m.sql); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed executing migration %d (%s): %w", m.version, m.name, err)
		}

		recordSQL := "INSERT INTO schema_migrations (version, name, applied_at) VALUES (?, ?, ?)"
		if _, err := tx.ExecContext(ctx, recordSQL, m.version, m.name, time.Now().UTC()); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed recording migration %d: %w", m.version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed committing migration %d: %w", m.version, err)
		}
	}

	return nil
}
