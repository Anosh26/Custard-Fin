package prefetch_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/adapters/memory"
	"github.com/anosh26/custardfin/apps/api/internal/adapters/prefetch"
)

type mockTunnelPuller struct {
	pullCount   int64
	simulateErr error
	payload     []byte
	contentType string
}

func (m *mockTunnelPuller) PullChunk(ctx context.Context, nodeID, mediaID, chunkPath string) (io.ReadCloser, string, int64, error) {
	atomic.AddInt64(&m.pullCount, 1)
	if m.simulateErr != nil {
		return nil, "", 0, m.simulateErr
	}
	data := m.payload
	if len(data) == 0 {
		data = []byte("dummy-ts-segment-bytes")
	}
	ct := m.contentType
	if ct == "" {
		ct = "video/mp2t"
	}
	return io.NopCloser(bytes.NewReader(data)), ct, int64(len(data)), nil
}

func TestNewBackgroundPrefetcher_Validation(t *testing.T) {
	cache := memory.NewMemoryChunkCache(time.Hour)
	puller := &mockTunnelPuller{}

	if _, err := prefetch.NewBackgroundPrefetcher(prefetch.Config{}, nil, puller, nil); err == nil {
		t.Fatalf("expected error for nil cache")
	}
	if _, err := prefetch.NewBackgroundPrefetcher(prefetch.Config{}, cache, nil, nil); err == nil {
		t.Fatalf("expected error for nil puller")
	}
}

func TestBackgroundPrefetcher_Prefetch(t *testing.T) {
	ctx := context.Background()

	t.Run("successfully prefetches uncached chunk", func(t *testing.T) {
		cache := memory.NewMemoryChunkCache(time.Hour)
		puller := &mockTunnelPuller{payload: []byte("segment-0-data")}

		p, err := prefetch.NewBackgroundPrefetcher(prefetch.Config{Timeout: 2 * time.Second}, cache, puller, nil)
		if err != nil {
			t.Fatalf("failed constructing prefetcher: %v", err)
		}

		err = p.Prefetch(ctx, "node-1", "media-1", "1080p/segment_000.ts")
		if err != nil {
			t.Fatalf("unexpected error during prefetch: %v", err)
		}

		if atomic.LoadInt64(&puller.pullCount) != 1 {
			t.Errorf("pullCount = %d, want 1", puller.pullCount)
		}

		// Verify chunk exists in cache
		exists, err := cache.Exists(ctx, "node-1/media-1/1080p/segment_000.ts")
		if err != nil || !exists {
			t.Errorf("expected chunk to exist in cache, err: %v", err)
		}

		// Subsequent prefetch should be skipped (already in cache)
		err = p.Prefetch(ctx, "node-1", "media-1", "1080p/segment_000.ts")
		if err != nil {
			t.Fatalf("unexpected error during second prefetch: %v", err)
		}
		if atomic.LoadInt64(&puller.pullCount) != 1 {
			t.Errorf("expected pullCount to remain 1 on cache hit, got %d", puller.pullCount)
		}
	})

	t.Run("handles invalid path safely", func(t *testing.T) {
		cache := memory.NewMemoryChunkCache(time.Hour)
		puller := &mockTunnelPuller{}
		p, _ := prefetch.NewBackgroundPrefetcher(prefetch.Config{}, cache, puller, nil)

		err := p.Prefetch(ctx, "node-1", "media-1", "../secret.txt")
		if err == nil {
			t.Fatalf("expected error for path traversal, got nil")
		}
	})

	t.Run("propagates tunnel errors gracefully", func(t *testing.T) {
		cache := memory.NewMemoryChunkCache(time.Hour)
		puller := &mockTunnelPuller{simulateErr: errors.New("tunnel connection dropped")}
		p, _ := prefetch.NewBackgroundPrefetcher(prefetch.Config{}, cache, puller, nil)

		err := p.Prefetch(ctx, "node-1", "media-1", "1080p/segment_000.ts")
		if err == nil {
			t.Fatalf("expected error on tunnel failure, got nil")
		}
	})
}

func TestBackgroundPrefetcher_TriggerPrefetch(t *testing.T) {
	ctx := context.Background()
	cache := memory.NewMemoryChunkCache(time.Hour)
	puller := &mockTunnelPuller{payload: []byte("async-chunk-data")}

	p, err := prefetch.NewBackgroundPrefetcher(prefetch.Config{Timeout: 2 * time.Second}, cache, puller, nil)
	if err != nil {
		t.Fatalf("failed constructing prefetcher: %v", err)
	}

	p.TriggerPrefetch(ctx, "node-async", "media-async", "720p/segment_000.ts")
	p.Wait()

	exists, err := cache.Exists(ctx, "node-async/media-async/720p/segment_000.ts")
	if err != nil || !exists {
		t.Errorf("expected async prefetched chunk to exist in cache, err: %v", err)
	}
}
