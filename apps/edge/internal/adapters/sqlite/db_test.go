package sqlite

import (
	"context"
	"path/filepath"
	"testing"
)

func TestOpen_InMemory(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite db: %v", err)
	}
	defer db.Close()

	// Verify foreign keys pragma
	var fk int
	if err := db.QueryRowContext(ctx, "PRAGMA foreign_keys;").Scan(&fk); err != nil {
		t.Fatalf("failed to query foreign_keys pragma: %v", err)
	}
	if fk != 1 {
		t.Errorf("expected foreign_keys = 1, got %d", fk)
	}

	// Verify busy_timeout pragma
	var timeout int
	if err := db.QueryRowContext(ctx, "PRAGMA busy_timeout;").Scan(&timeout); err != nil {
		t.Fatalf("failed to query busy_timeout pragma: %v", err)
	}
	if timeout < 5000 {
		t.Errorf("expected busy_timeout >= 5000, got %d", timeout)
	}
}

func TestOpen_OnDiskWithAutoDirCreation(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "nested", "folder", "test_edge.db")

	ctx := context.Background()
	db, err := Open(ctx, dbPath)
	if err != nil {
		t.Fatalf("failed to open on-disk sqlite db: %v", err)
	}
	defer db.Close()

	// Verify WAL journal mode on disk
	var journalMode string
	if err := db.QueryRowContext(ctx, "PRAGMA journal_mode;").Scan(&journalMode); err != nil {
		t.Fatalf("failed to query journal_mode pragma: %v", err)
	}
	if journalMode != "wal" {
		t.Errorf("expected journal_mode 'wal', got %q", journalMode)
	}
}

func TestMigrate_Idempotency(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// Re-run Migrate: must be idempotent and succeed without error
	if err := db.Migrate(ctx); err != nil {
		t.Fatalf("re-running Migrate failed: %v", err)
	}

	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations;").Scan(&count); err != nil {
		t.Fatalf("failed to query schema_migrations count: %v", err)
	}
	if count != len(migrations) {
		t.Errorf("expected %d migrations recorded, got %d", len(migrations), count)
	}
}
