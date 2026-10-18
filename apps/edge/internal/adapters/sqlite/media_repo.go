package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/domain"
	"github.com/anosh26/custardfin/apps/edge/internal/ports"
)

var _ ports.MediaRepository = (*MediaRepository)(nil)

// MediaRepository provides SQLite persistence for media items.
type MediaRepository struct {
	db *DB
}

// NewMediaRepository constructs a MediaRepository instance.
func NewMediaRepository(db *DB) *MediaRepository {
	return &MediaRepository{db: db}
}

// Save inserts or updates a media entity.
func (r *MediaRepository) Save(ctx context.Context, m *domain.Media) error {
	query := `
	INSERT INTO media (
		id, path, file_size, mod_time, content_hash, title, duration_s,
		width, height, video_codec, audio_codec, is_public, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		path = excluded.path,
		file_size = excluded.file_size,
		mod_time = excluded.mod_time,
		content_hash = excluded.content_hash,
		title = excluded.title,
		duration_s = excluded.duration_s,
		width = excluded.width,
		height = excluded.height,
		video_codec = excluded.video_codec,
		audio_codec = excluded.audio_codec,
		is_public = excluded.is_public,
		updated_at = excluded.updated_at;
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		m.ID,
		m.Path,
		m.FileSize,
		m.ModTime.Format(time.RFC3339Nano),
		m.ContentHash,
		m.Title,
		m.DurationS,
		m.Width,
		m.Height,
		m.VideoCodec,
		m.AudioCodec,
		false,
		m.CreatedAt.Format(time.RFC3339Nano),
		m.UpdatedAt.Format(time.RFC3339Nano),
	)
	if err != nil {
		return fmt.Errorf("failed to save media %s: %w", m.ID, err)
	}
	return nil
}

// FindByPath retrieves a media entity by its local absolute filesystem path.
func (r *MediaRepository) FindByPath(ctx context.Context, path string) (*domain.Media, error) {
	cleanPath := filepath.Clean(path)
	query := `
	SELECT id, path, file_size, mod_time, content_hash, title, duration_s,
	       width, height, video_codec, audio_codec, created_at, updated_at
	FROM media
	WHERE path = ?;
	`
	return r.scanSingle(r.db.QueryRowContext(ctx, query, cleanPath))
}

// FindByID retrieves a media entity by its UUID.
func (r *MediaRepository) FindByID(ctx context.Context, id string) (*domain.Media, error) {
	query := `
	SELECT id, path, file_size, mod_time, content_hash, title, duration_s,
	       width, height, video_codec, audio_codec, created_at, updated_at
	FROM media
	WHERE id = ?;
	`
	return r.scanSingle(r.db.QueryRowContext(ctx, query, id))
}

// FindByHash retrieves a media entity by its SHA-256 content hash.
func (r *MediaRepository) FindByHash(ctx context.Context, hash string) (*domain.Media, error) {
	query := `
	SELECT id, path, file_size, mod_time, content_hash, title, duration_s,
	       width, height, video_codec, audio_codec, created_at, updated_at
	FROM media
	WHERE content_hash = ?;
	`
	return r.scanSingle(r.db.QueryRowContext(ctx, query, hash))
}

// List returns all cataloged media entities ordered by updated timestamp descending.
func (r *MediaRepository) List(ctx context.Context) ([]*domain.Media, error) {
	query := `
	SELECT id, path, file_size, mod_time, content_hash, title, duration_s,
	       width, height, video_codec, audio_codec, created_at, updated_at
	FROM media
	ORDER BY updated_at DESC;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed querying media list: %w", err)
	}
	defer rows.Close()

	var list []*domain.Media
	for rows.Next() {
		m, err := scanMedia(rows)
		if err != nil {
			return nil, fmt.Errorf("failed scanning media row: %w", err)
		}
		list = append(list, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return list, nil
}

// Delete removes a media entity by UUID.
func (r *MediaRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM media WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed deleting media %s: %w", id, err)
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func (r *MediaRepository) scanSingle(row *sql.Row) (*domain.Media, error) {
	m, err := scanMedia(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan media: %w", err)
	}
	return m, nil
}

func scanMedia(s rowScanner) (*domain.Media, error) {
	var (
		m          domain.Media
		modTimeStr string
		createStr  string
		updateStr  string
	)

	err := s.Scan(
		&m.ID,
		&m.Path,
		&m.FileSize,
		&modTimeStr,
		&m.ContentHash,
		&m.Title,
		&m.DurationS,
		&m.Width,
		&m.Height,
		&m.VideoCodec,
		&m.AudioCodec,
		&createStr,
		&updateStr,
	)
	if err != nil {
		return nil, err
	}

	m.ModTime, _ = time.Parse(time.RFC3339Nano, modTimeStr)
	m.CreatedAt, _ = time.Parse(time.RFC3339Nano, createStr)
	m.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updateStr)

	return &m, nil
}
