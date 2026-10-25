package ports

import (
	"context"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
)

// SaturationEstimator assesses the concurrent load and capacity headroom of a home edge node.
type SaturationEstimator interface {
	// EstimateLoad evaluates current node liveness, stream saturation, and returns its LoadLevel and active stream count.
	EstimateLoad(ctx context.Context, nodeID string) (domain.LoadLevel, int, error)
}

// Prefetcher asynchronously preheats hot initial media segments into the Origin Shield cache to eliminate startup buffering.
type Prefetcher interface {
	// TriggerPrefetch queues an asynchronous prefetch for the first media chunk into the cache.
	TriggerPrefetch(ctx context.Context, nodeID, mediaID, initialChunkPath string)
}
