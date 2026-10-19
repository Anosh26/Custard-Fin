package memory

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
)

func TestMemoryChunkCache_BasicLifecycle(t *testing.T) {
	cache := NewMemoryChunkCache(24 * time.Hour)
	ctx := context.Background()

	key := "node1/media1/segment_0000.ts"
	payload := []byte("fake-ts-binary-data")

	// 1. Exists before Put should be false
	exists, err := cache.Exists(ctx, key)
	if err != nil {
		t.Fatalf("unexpected error on Exists: %v", err)
	}
	if exists {
		t.Fatalf("expected key to not exist")
	}

	// 2. Get before Put should return domain.ErrNotFound
	_, err = cache.Get(ctx, key)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// 3. Put chunk
	err = cache.Put(ctx, key, bytes.NewReader(payload), int64(len(payload)), "video/mp2t")
	if err != nil {
		t.Fatalf("unexpected error on Put: %v", err)
	}

	// 4. Exists after Put should be true
	exists, err = cache.Exists(ctx, key)
	if err != nil || !exists {
		t.Fatalf("expected key to exist after Put, err: %v", err)
	}

	// 5. Get chunk and verify content & headers
	item, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("unexpected error on Get: %v", err)
	}
	defer item.Body.Close()

	if item.ContentType != "video/mp2t" {
		t.Errorf("expected ContentType video/mp2t, got %s", item.ContentType)
	}
	if item.ContentLength != int64(len(payload)) {
		t.Errorf("expected ContentLength %d, got %d", len(payload), item.ContentLength)
	}

	readData, err := io.ReadAll(item.Body)
	if err != nil {
		t.Fatalf("failed reading body: %v", err)
	}
	if !bytes.Equal(readData, payload) {
		t.Errorf("read data %s does not match expected payload %s", readData, payload)
	}

	// 6. Delete chunk
	if err := cache.Delete(ctx, key); err != nil {
		t.Fatalf("unexpected error on Delete: %v", err)
	}

	exists, err = cache.Exists(ctx, key)
	if err != nil || exists {
		t.Fatalf("expected key to be removed after Delete")
	}
}

func TestMemoryChunkCache_Expiration(t *testing.T) {
	// 50ms short TTL to test auto-expiration
	cache := NewMemoryChunkCache(50 * time.Millisecond)
	ctx := context.Background()

	key := "node1/media1/chunk.ts"
	payload := []byte("short-lived-data")

	if err := cache.Put(ctx, key, bytes.NewReader(payload), int64(len(payload)), "video/mp2t"); err != nil {
		t.Fatalf("failed to Put: %v", err)
	}

	// Immediate check: exists
	exists, err := cache.Exists(ctx, key)
	if err != nil || !exists {
		t.Fatalf("expected item to exist before expiry")
	}

	// Wait for TTL to pass
	time.Sleep(70 * time.Millisecond)

	// Exists should now be false
	exists, err = cache.Exists(ctx, key)
	if err != nil || exists {
		t.Fatalf("expected item to be expired")
	}

	// Get should return ErrNotFound
	_, err = cache.Get(ctx, key)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for expired item, got %v", err)
	}
}

func TestMemoryChunkCache_ValidationAndErrors(t *testing.T) {
	cache := NewMemoryChunkCache(0) // Should default to 24h
	ctx := context.Background()

	// Nil reader should error
	err := cache.Put(ctx, "k1", nil, 0, "")
	if err == nil {
		t.Fatalf("expected error on nil reader")
	}

	// Simulated error
	simulated := errors.New("simulated memory error")
	cache.SetSimulatedError(simulated)

	_, err = cache.Exists(ctx, "k1")
	if !errors.Is(err, simulated) {
		t.Fatalf("expected simulated error on Exists, got %v", err)
	}

	_, err = cache.Get(ctx, "k1")
	if !errors.Is(err, simulated) {
		t.Fatalf("expected simulated error on Get, got %v", err)
	}

	err = cache.Put(ctx, "k1", bytes.NewReader([]byte("test")), 4, "text/plain")
	if !errors.Is(err, simulated) {
		t.Fatalf("expected simulated error on Put, got %v", err)
	}

	err = cache.Delete(ctx, "k1")
	if !errors.Is(err, simulated) {
		t.Fatalf("expected simulated error on Delete, got %v", err)
	}

	// Clear simulated error
	cache.SetSimulatedError(nil)

	// Context cancellation
	cancCtx, cancel := context.WithCancel(ctx)
	cancel()

	_, err = cache.Get(cancCtx, "k1")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestMemoryChunkCache_Concurrency(t *testing.T) {
	cache := NewMemoryChunkCache(1 * time.Hour)
	ctx := context.Background()

	var wg sync.WaitGroup
	workers := 20
	iterations := 50

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(wID int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				key := "worker-key"
				_ = cache.Put(ctx, key, bytes.NewReader([]byte("data")), 4, "video/mp2t")
				_, _ = cache.Exists(ctx, key)
				item, err := cache.Get(ctx, key)
				if err == nil && item != nil {
					_ = item.Body.Close()
				}
			}
		}(i)
	}

	wg.Wait()
	if cache.Len() == 0 {
		t.Fatalf("expected active items after concurrent writes")
	}
}
