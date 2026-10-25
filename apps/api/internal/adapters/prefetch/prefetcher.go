package prefetch

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/ports"
)

var (
	ErrNilDependency = errors.New("chunkCache and tunnelPuller are required for prefetcher")
)

// Config configures the background chunk prefetcher.
type Config struct {
	Timeout time.Duration // default 5s
}

// BackgroundPrefetcher proactively warms the Origin Shield cache with initial video chunks.
type BackgroundPrefetcher struct {
	cfg          Config
	chunkCache   ports.ChunkCache
	tunnelPuller ports.TunnelPuller
	logger       *slog.Logger
	wg           sync.WaitGroup
}

var _ ports.Prefetcher = (*BackgroundPrefetcher)(nil)

// NewBackgroundPrefetcher constructs a new BackgroundPrefetcher.
func NewBackgroundPrefetcher(cfg Config, chunkCache ports.ChunkCache, tunnelPuller ports.TunnelPuller, logger *slog.Logger) (*BackgroundPrefetcher, error) {
	if chunkCache == nil || tunnelPuller == nil {
		return nil, ErrNilDependency
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5 * time.Second
	}
	if logger == nil {
		logger = slog.Default()
	}

	return &BackgroundPrefetcher{
		cfg:          cfg,
		chunkCache:   chunkCache,
		tunnelPuller: tunnelPuller,
		logger:       logger,
	}, nil
}

// Prefetch synchronously pulls and caches the target chunk if not already present.
func (p *BackgroundPrefetcher) Prefetch(ctx context.Context, nodeID, mediaID, chunkPath string) error {
	cleanPath := path.Clean(strings.TrimPrefix(chunkPath, "/"))
	if cleanPath == "." || strings.HasPrefix(cleanPath, "..") {
		return fmt.Errorf("invalid prefetch path: %s", chunkPath)
	}

	cacheKey := fmt.Sprintf("%s/%s/%s", nodeID, mediaID, cleanPath)

	// Check if already in Origin Shield cache
	exists, err := p.chunkCache.Exists(ctx, cacheKey)
	if err == nil && exists {
		p.logger.Debug("prefetch skipped: chunk already cached", "key", cacheKey)
		return nil
	}

	// Pull chunk via tunnel puller
	body, contentType, _, err := p.tunnelPuller.PullChunk(ctx, nodeID, mediaID, cleanPath)
	if err != nil {
		return fmt.Errorf("tunnel pull failed during prefetch: %w", err)
	}
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed reading chunk body during prefetch: %w", err)
	}

	// Store into cache with 24h retention
	if err := p.chunkCache.Put(ctx, cacheKey, bytes.NewReader(data), int64(len(data)), contentType); err != nil {
		return fmt.Errorf("failed storing prefetched chunk in cache: %w", err)
	}

	p.logger.Info("preheated initial chunk into origin shield cache", "key", cacheKey, "bytes", len(data))
	return nil
}

// TriggerPrefetch initiates a non-blocking asynchronous prefetch in a detached background goroutine.
func (p *BackgroundPrefetcher) TriggerPrefetch(ctx context.Context, nodeID, mediaID, chunkPath string) {
	if strings.TrimSpace(nodeID) == "" || strings.TrimSpace(mediaID) == "" || strings.TrimSpace(chunkPath) == "" {
		return
	}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		prefetchCtx, cancel := context.WithTimeout(context.Background(), p.cfg.Timeout)
		defer cancel()

		if err := p.Prefetch(prefetchCtx, nodeID, mediaID, chunkPath); err != nil {
			p.logger.Debug("non-fatal prefetch failure", "node_id", nodeID, "media_id", mediaID, "chunk_path", chunkPath, "err", err)
		}
	}()
}

// Wait blocks until all active background prefetch routines finish (useful for graceful shutdown and tests).
func (p *BackgroundPrefetcher) Wait() {
	p.wg.Wait()
}
