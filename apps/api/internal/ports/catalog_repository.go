package ports

import (
	"context"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
)

// CatalogRepository provides persistence operations for public media catalog manifests and renditions.
type CatalogRepository interface {
	UpsertManifest(ctx context.Context, item *domain.CatalogItem) error
	DeleteManifest(ctx context.Context, nodeID, mediaID string) error
	Search(ctx context.Context, query string, limit, offset int) ([]*domain.CatalogItem, int, error)
	FindByMediaID(ctx context.Context, mediaID string) (*domain.CatalogItem, error)
}
