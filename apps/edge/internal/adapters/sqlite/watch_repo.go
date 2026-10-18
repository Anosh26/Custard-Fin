package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/domain"
	"github.com/anosh26/custardfin/apps/edge/internal/ports"
)

var _ ports.WatchProgressRepository = (*WatchProgressRepository)(nil)

// WatchProgressRepository stores private watch progress in edge SQLite.
type WatchProgressRepository struct {
	db *DB
}

// NewWatchProgressRepository constructs a WatchProgressRepository instance.
func NewWatchProgressRepository(db *DB) *WatchProgressRepository {
	return &WatchProgressRepository{db: db}
}

// GetWatchProgress retrieves playback progress for a given media item.
func (r *WatchProgressRepository) GetWatchProgress(ctx context.Context, mediaID string) (*domain.WatchProgress, error) {
	query := `
	SELECT media_id, position_s, completed, last_watched_at
	FROM watch_progress
	WHERE media_id = ?;
	`

	var (
		wp         domain.WatchProgress
		watchedStr string
	)

	err := r.db.QueryRowContext(ctx, query, mediaID).Scan(
		&wp.MediaID,
		&wp.PositionS,
		&wp.Completed,
		&watchedStr,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query watch progress: %w", err)
	}

	wp.LastWatchedAt, _ = time.Parse(time.RFC3339Nano, watchedStr)
	return &wp, nil
}

// SaveWatchProgress inserts or updates playback progress.
func (r *WatchProgressRepository) SaveWatchProgress(ctx context.Context, wp *domain.WatchProgress) error {
	query := `
	INSERT INTO watch_progress (media_id, position_s, completed, last_watched_at)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(media_id) DO UPDATE SET
		position_s = excluded.position_s,
		completed = excluded.completed,
		last_watched_at = excluded.last_watched_at;
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		wp.MediaID,
		wp.PositionS,
		wp.Completed,
		wp.LastWatchedAt.Format(time.RFC3339Nano),
	)
	if err != nil {
		return fmt.Errorf("failed to save watch progress for media %s: %w", wp.MediaID, err)
	}
	return nil
}

// DeleteWatchProgress removes playback progress for a given media item.
func (r *WatchProgressRepository) DeleteWatchProgress(ctx context.Context, mediaID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM watch_progress WHERE media_id = ?", mediaID)
	if err != nil {
		return fmt.Errorf("failed deleting watch progress for media %s: %w", mediaID, err)
	}
	return nil
}
