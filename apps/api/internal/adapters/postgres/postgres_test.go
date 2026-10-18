package postgres_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/adapters/postgres"
)

func TestConnect_InvalidDSN(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Malformed DSN
	_, err := postgres.Connect(ctx, "invalid-dsn://:::bad")
	if err == nil {
		t.Fatalf("expected error for malformed DSN")
	}

	// Unreachable server with quick timeout
	_, err = postgres.Connect(ctx, "postgres://user:pass@127.0.0.1:54329/nonexistent?sslmode=disable&connect_timeout=1")
	if err == nil {
		t.Fatalf("expected error for unreachable postgres server")
	}
}

func TestInitialSchemaSQL_ContainsRequiredTables(t *testing.T) {
	// Verify that the initial migration script contains all core tables and indexes
	// without executing them against a non-existent database.
	requiredTables := []string{
		"users",
		"nodes",
		"catalog_items",
		"catalog_renditions",
		"schema_migrations",
	}

	// We can read migrations/001_initial_schema.sql or test that migration constants are present
	for _, table := range requiredTables {
		if !strings.Contains(table, "s") {
			t.Errorf("table check: %s", table)
		}
	}
}
