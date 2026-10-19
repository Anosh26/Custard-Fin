package memory

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
)

type memoryCacheEntry struct {
	data        []byte
	contentType string
	expiresAt   time.Time
}

// ChunkCache is an in-memory, thread-safe implementation of ports.ChunkCache.
// Adheres to Invariant 1 with default 24h TTL.
type ChunkCache struct {
	mu           sync.RWMutex
	items        map[string]memoryCacheEntry
	ttl          time.Duration
	simulatedErr error
}

// NewMemoryChunkCache creates a new in-memory chunk cache.
func NewMemoryChunkCache(ttl time.Duration) *ChunkCache {
	if ttl <= 0 {
		ttl = 24 * time.Hour // Invariant 1: 24h default TTL
	}
	return &ChunkCache{
		items: make(map[string]memoryCacheEntry),
		ttl:   ttl,
	}
}

// SetSimulatedError configures an error to return on all operations for test simulation.
func (c *ChunkCache) SetSimulatedError(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.simulatedErr = err
}

// Get retrieves a chunk by key.
func (c *ChunkCache) Get(ctx context.Context, key string) (*ports.CacheItem, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.simulatedErr != nil {
		return nil, c.simulatedErr
	}

	entry, ok := c.items[key]
	if !ok {
		return nil, domain.ErrNotFound
	}

	if time.Now().After(entry.expiresAt) {
		delete(c.items, key)
		return nil, domain.ErrNotFound
	}

	return &ports.CacheItem{
		Body:          io.NopCloser(bytes.NewReader(entry.data)),
		ContentType:   entry.contentType,
		ContentLength: int64(len(entry.data)),
	}, nil
}

// Put stores a chunk in the cache with the configured TTL (24h).
func (c *ChunkCache) Put(ctx context.Context, key string, r io.Reader, contentLength int64, contentType string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if r == nil {
		return errors.New("cannot put nil reader into cache")
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.simulatedErr != nil {
		return c.simulatedErr
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.items[key] = memoryCacheEntry{
		data:        data,
		contentType: contentType,
		expiresAt:   time.Now().Add(c.ttl),
	}
	return nil
}

// Exists checks if a chunk is cached and not expired.
func (c *ChunkCache) Exists(ctx context.Context, key string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.simulatedErr != nil {
		return false, c.simulatedErr
	}

	entry, ok := c.items[key]
	if !ok {
		return false, nil
	}

	if time.Now().After(entry.expiresAt) {
		delete(c.items, key)
		return false, nil
	}

	return true, nil
}

// Delete removes a chunk from the cache.
func (c *ChunkCache) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.simulatedErr != nil {
		return c.simulatedErr
	}

	delete(c.items, key)
	return nil
}

// Len returns the count of active items in the cache.
func (c *ChunkCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
