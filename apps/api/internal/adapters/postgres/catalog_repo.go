package postgres

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogRepository struct {
	pool *pgxpool.Pool
}

func NewCatalogRepository(pool *pgxpool.Pool) ports.CatalogRepository {
	return &CatalogRepository{pool: pool}
}

func (r *CatalogRepository) UpsertManifest(ctx context.Context, item *domain.CatalogItem) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin upsert manifest tx: %w", err)
	}
	defer tx.Rollback(ctx)

	upsertItemSQL := `
	INSERT INTO catalog_items (
		media_id, node_id, title, description, duration_s, content_hash, is_public, manifest_json, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	ON CONFLICT (node_id, media_id) DO UPDATE SET
		title = EXCLUDED.title,
		description = EXCLUDED.description,
		duration_s = EXCLUDED.duration_s,
		content_hash = EXCLUDED.content_hash,
		is_public = EXCLUDED.is_public,
		manifest_json = EXCLUDED.manifest_json,
		updated_at = EXCLUDED.updated_at
	`
	_, err = tx.Exec(ctx, upsertItemSQL,
		item.MediaID,
		item.NodeID,
		item.Title,
		item.Description,
		item.DurationS,
		item.ContentHash,
		item.IsPublic,
		item.ManifestJSON,
		item.CreatedAt,
		item.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert catalog item: %w", err)
	}

	// Remove old renditions
	deleteRenditionsSQL := `DELETE FROM catalog_renditions WHERE node_id = $1 AND media_id = $2`
	if _, err := tx.Exec(ctx, deleteRenditionsSQL, item.NodeID, item.MediaID); err != nil {
		return fmt.Errorf("failed to clear old renditions: %w", err)
	}

	// Insert new renditions
	insertRenditionSQL := `
	INSERT INTO catalog_renditions (
		id, media_id, node_id, resolution, bitrate_bps, playlist_path, chunk_duration_s, created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	for _, rend := range item.Renditions {
		rendID := rendUUID()
		_, err = tx.Exec(ctx, insertRenditionSQL,
			rendID,
			item.MediaID,
			item.NodeID,
			rend.Resolution,
			rend.BitrateBps,
			rend.PlaylistPath,
			rend.ChunkDurationS,
			item.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert rendition: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *CatalogRepository) DeleteManifest(ctx context.Context, nodeID, mediaID string) error {
	query := `DELETE FROM catalog_items WHERE node_id = $1 AND media_id = $2`
	tag, err := r.pool.Exec(ctx, query, nodeID, mediaID)
	if err != nil {
		return fmt.Errorf("failed to delete catalog item: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CatalogRepository) Search(ctx context.Context, query string, limit, offset int, allowedNodeIDs []string) ([]*domain.CatalogItem, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	query = strings.TrimSpace(query)
	var countSQL string
	var searchSQL string
	var countArgs []interface{}
	var searchArgs []interface{}

	if len(allowedNodeIDs) == 0 {
		if query == "" {
			countSQL = `SELECT COUNT(*) FROM catalog_items WHERE is_public = TRUE`
			searchSQL = `
			SELECT media_id, node_id, title, description, duration_s, content_hash, is_public, manifest_json, created_at, updated_at
			FROM catalog_items
			WHERE is_public = TRUE
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
			`
			searchArgs = []interface{}{limit, offset}
		} else {
			likePattern := "%" + query + "%"
			countSQL = `SELECT COUNT(*) FROM catalog_items WHERE is_public = TRUE AND (title ILIKE $1 OR description ILIKE $1)`
			countArgs = []interface{}{likePattern}

			searchSQL = `
			SELECT media_id, node_id, title, description, duration_s, content_hash, is_public, manifest_json, created_at, updated_at
			FROM catalog_items
			WHERE is_public = TRUE AND (title ILIKE $1 OR description ILIKE $1)
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
			`
			searchArgs = []interface{}{likePattern, limit, offset}
		}
	} else {
		if query == "" {
			countSQL = `SELECT COUNT(*) FROM catalog_items WHERE is_public = TRUE OR node_id = ANY($1)`
			countArgs = []interface{}{allowedNodeIDs}

			searchSQL = `
			SELECT media_id, node_id, title, description, duration_s, content_hash, is_public, manifest_json, created_at, updated_at
			FROM catalog_items
			WHERE is_public = TRUE OR node_id = ANY($1)
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
			`
			searchArgs = []interface{}{allowedNodeIDs, limit, offset}
		} else {
			likePattern := "%" + query + "%"
			countSQL = `SELECT COUNT(*) FROM catalog_items WHERE (is_public = TRUE OR node_id = ANY($2)) AND (title ILIKE $1 OR description ILIKE $1)`
			countArgs = []interface{}{likePattern, allowedNodeIDs}

			searchSQL = `
			SELECT media_id, node_id, title, description, duration_s, content_hash, is_public, manifest_json, created_at, updated_at
			FROM catalog_items
			WHERE (is_public = TRUE OR node_id = ANY($2)) AND (title ILIKE $1 OR description ILIKE $1)
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4
			`
			searchArgs = []interface{}{likePattern, allowedNodeIDs, limit, offset}
		}
	}

	var total int
	if err := r.pool.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count catalog items: %w", err)
	}

	rows, err := r.pool.Query(ctx, searchSQL, searchArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search catalog items: %w", err)
	}
	defer rows.Close()

	var items []*domain.CatalogItem
	for rows.Next() {
		var it domain.CatalogItem
		if err := rows.Scan(
			&it.MediaID,
			&it.NodeID,
			&it.Title,
			&it.Description,
			&it.DurationS,
			&it.ContentHash,
			&it.IsPublic,
			&it.ManifestJSON,
			&it.CreatedAt,
			&it.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan catalog item: %w", err)
		}
		items = append(items, &it)
	}

	return items, total, rows.Err()
}

func (r *CatalogRepository) FindByMediaID(ctx context.Context, mediaID string) (*domain.CatalogItem, error) {
	mediaID = strings.TrimSpace(mediaID)
	query := `
	SELECT media_id, node_id, title, description, duration_s, content_hash, is_public, manifest_json, created_at, updated_at
	FROM catalog_items
	WHERE media_id = $1
	`
	row := r.pool.QueryRow(ctx, query, mediaID)

	var it domain.CatalogItem
	err := row.Scan(
		&it.MediaID,
		&it.NodeID,
		&it.Title,
		&it.Description,
		&it.DurationS,
		&it.ContentHash,
		&it.IsPublic,
		&it.ManifestJSON,
		&it.CreatedAt,
		&it.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to query catalog item: %w", err)
	}

	// Fetch renditions
	rendSQL := `
	SELECT resolution, bitrate_bps, playlist_path, chunk_duration_s
	FROM catalog_renditions
	WHERE node_id = $1 AND media_id = $2
	ORDER BY bitrate_bps DESC
	`
	rows, err := r.pool.Query(ctx, rendSQL, it.NodeID, it.MediaID)
	if err != nil {
		return nil, fmt.Errorf("failed to query renditions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rend domain.Rendition
		if err := rows.Scan(&rend.Resolution, &rend.BitrateBps, &rend.PlaylistPath, &rend.ChunkDurationS); err != nil {
			return nil, fmt.Errorf("failed to scan rendition: %w", err)
		}
		it.Renditions = append(it.Renditions, rend)
	}

	return &it, rows.Err()
}

func rendUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
