# apps/edge/internal/adapters/sqlite

Pure-Go SQLite database store for CustardFin Edge daemon using `modernc.org/sqlite`.

## Purpose
Maintains local metadata on discovered media items and private user watch history.
Adheres strictly to **Invariant 8: Private data stays home** (watch progress never syncs to the central cloud).

## Public Interface
- `Open(ctx context.Context, dsn string) (*DB, error)`: Configures connection with `PRAGMA journal_mode=WAL;`, `PRAGMA busy_timeout=5000;`, and `PRAGMA foreign_keys=ON;`.
- `DB.Migrate(ctx context.Context) error`: Applies versioned idempotent schema migrations.
- `NewMediaRepository(db *DB) *MediaRepository`: Implements `ports.MediaRepository`.
- `NewWatchProgressRepository(db *DB) *WatchProgressRepository`: Implements `ports.WatchProgressRepository`.

## Testing
```bash
go test -v ./apps/edge/internal/adapters/sqlite
```
