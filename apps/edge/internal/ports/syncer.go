package ports

import (
	"context"

	"github.com/anosh26/custardfin/apps/edge/internal/domain"
	"github.com/anosh26/custardfin/libs/contracts"
)

// ManifestSyncer provides syncing operations between the local edge daemon and the Cloud Control API.
type ManifestSyncer interface {
	SyncManifest(ctx context.Context, media *domain.Media, renditions []contracts.Rendition) error
	UnsyncManifest(ctx context.Context, mediaID string) error
}
