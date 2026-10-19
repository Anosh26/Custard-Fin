package ports

import (
	"context"
	"io"
)

// CacheItem represents a cached video segment or playlist item.
type CacheItem struct {
	Body          io.ReadCloser
	ContentType   string
	ContentLength int64
}

// ChunkCache defines the contract for temporary HLS chunk storage in the Origin Shield.
// Implementations MUST adhere to Invariant 1 (temporary cache only, 24h retention)
// and Invariant 7 (hot chunks stored once to shield home node upload).
type ChunkCache interface {
	// Get retrieves a cached chunk by key. If the key is not cached, returns domain.ErrNotFound.
	Get(ctx context.Context, key string) (*CacheItem, error)

	// Put stores a chunk in the cache with the given content type.
	// Storage backends tag objects with a 24-hour expiration policy.
	Put(ctx context.Context, key string, r io.Reader, contentLength int64, contentType string) error

	// Exists returns true if the key is already stored in the cache.
	Exists(ctx context.Context, key string) (bool, error)

	// Delete removes a chunk from the cache.
	Delete(ctx context.Context, key string) error
}

// TunnelPuller defines the contract for pulling an uncached chunk from a home edge node
// through the FRP reverse tunnel relay (Invariant 2, Invariant 7).
type TunnelPuller interface {
	// PullChunk requests a media chunk or playlist from the edge node via the tunnel.
	// Returns an io.ReadCloser of the stream, content type, content length, and error.
	PullChunk(ctx context.Context, nodeID, mediaID, chunkPath string) (io.ReadCloser, string, int64, error)
}
