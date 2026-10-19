package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
)

type mockS3Item struct {
	data        []byte
	contentType string
	tagging     string
}

type mockS3API struct {
	mu           sync.RWMutex
	items        map[string]mockS3Item
	simulatedErr error
}

func newMockS3API() *mockS3API {
	return &mockS3API{
		items: make(map[string]mockS3Item),
	}
}

func (m *mockS3API) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.simulatedErr != nil {
		return nil, m.simulatedErr
	}

	key := aws.ToString(params.Key)
	item, ok := m.items[key]
	if !ok {
		return nil, &types.NoSuchKey{Message: aws.String("The specified key does not exist.")}
	}

	return &s3.GetObjectOutput{
		Body:          io.NopCloser(bytes.NewReader(item.data)),
		ContentType:   aws.String(item.contentType),
		ContentLength: aws.Int64(int64(len(item.data))),
	}, nil
}

func (m *mockS3API) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.simulatedErr != nil {
		return nil, m.simulatedErr
	}

	data, err := io.ReadAll(params.Body)
	if err != nil {
		return nil, err
	}

	key := aws.ToString(params.Key)
	m.items[key] = mockS3Item{
		data:        data,
		contentType: aws.ToString(params.ContentType),
		tagging:     aws.ToString(params.Tagging),
	}

	return &s3.PutObjectOutput{}, nil
}

func (m *mockS3API) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.simulatedErr != nil {
		return nil, m.simulatedErr
	}

	key := aws.ToString(params.Key)
	item, ok := m.items[key]
	if !ok {
		return nil, &types.NotFound{Message: aws.String("Not Found")}
	}

	return &s3.HeadObjectOutput{
		ContentType:   aws.String(item.contentType),
		ContentLength: aws.Int64(int64(len(item.data))),
	}, nil
}

func (m *mockS3API) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.simulatedErr != nil {
		return nil, m.simulatedErr
	}

	delete(m.items, aws.ToString(params.Key))
	return &s3.DeleteObjectOutput{}, nil
}

func TestS3Config_Validation(t *testing.T) {
	// Missing bucket
	cfg := Config{}
	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected error on empty bucket")
	}

	// Valid with default region
	cfg = Config{Bucket: "test-bucket"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if cfg.Region != "us-east-1" {
		t.Errorf("expected default region us-east-1, got %s", cfg.Region)
	}

	// Constructor check
	cache, err := NewS3ChunkCache(cfg)
	if err != nil || cache == nil {
		t.Fatalf("expected successful cache constructor, err: %v", err)
	}
}

func TestS3ChunkCache_Lifecycle(t *testing.T) {
	mockAPI := newMockS3API()
	cache := NewS3ChunkCacheWithClient(mockAPI, "custardfin-cache")
	ctx := context.Background()

	key := "node-uuid/media-uuid/1080p/segment_0001.ts"
	data := []byte("ts-binary-segment-content")

	// 1. Initial Exists -> false
	exists, err := cache.Exists(ctx, key)
	if err != nil || exists {
		t.Fatalf("expected key to not exist, err: %v", err)
	}

	// 2. Initial Get -> ErrNotFound
	_, err = cache.Get(ctx, key)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// 3. Put chunk
	err = cache.Put(ctx, key, bytes.NewReader(data), int64(len(data)), "video/mp2t")
	if err != nil {
		t.Fatalf("failed to Put chunk: %v", err)
	}

	// Verify retention tagging was set
	item := mockAPI.items[key]
	if !strings.Contains(item.tagging, "retention=24h") {
		t.Errorf("expected retention=24h tagging, got %q", item.tagging)
	}

	// 4. Exists -> true
	exists, err = cache.Exists(ctx, key)
	if err != nil || !exists {
		t.Fatalf("expected key to exist after Put, err: %v", err)
	}

	// 5. Get chunk
	cachedItem, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("failed to Get chunk: %v", err)
	}
	defer cachedItem.Body.Close()

	if cachedItem.ContentType != "video/mp2t" {
		t.Errorf("expected ContentType video/mp2t, got %s", cachedItem.ContentType)
	}
	if cachedItem.ContentLength != int64(len(data)) {
		t.Errorf("expected ContentLength %d, got %d", len(data), cachedItem.ContentLength)
	}

	readData, err := io.ReadAll(cachedItem.Body)
	if err != nil || !bytes.Equal(readData, data) {
		t.Fatalf("mismatched data: %s vs %s", readData, data)
	}

	// 6. Delete chunk
	if err := cache.Delete(ctx, key); err != nil {
		t.Fatalf("failed to Delete: %v", err)
	}

	exists, err = cache.Exists(ctx, key)
	if err != nil || exists {
		t.Fatalf("expected key to be deleted")
	}
}

func TestS3ChunkCache_ValidationAndErrorHandling(t *testing.T) {
	mockAPI := newMockS3API()
	cache := NewS3ChunkCacheWithClient(mockAPI, "test-bucket")
	ctx := context.Background()

	// Empty key validation
	if _, err := cache.Get(ctx, ""); err == nil {
		t.Errorf("expected error on empty key for Get")
	}
	if err := cache.Put(ctx, "", bytes.NewReader([]byte("x")), 1, ""); err == nil {
		t.Errorf("expected error on empty key for Put")
	}
	if err := cache.Put(ctx, "valid-key", nil, 0, ""); err == nil {
		t.Errorf("expected error on nil reader for Put")
	}
	if _, err := cache.Exists(ctx, ""); err == nil {
		t.Errorf("expected error on empty key for Exists")
	}
	if err := cache.Delete(ctx, ""); err == nil {
		t.Errorf("expected error on empty key for Delete")
	}

	// Simulated S3 error
	simulated := errors.New("s3 connection reset")
	mockAPI.simulatedErr = simulated

	if _, err := cache.Get(ctx, "k"); !errors.Is(err, simulated) {
		t.Errorf("expected simulated error on Get, got %v", err)
	}
	if err := cache.Put(ctx, "k", bytes.NewReader([]byte("x")), 1, ""); !errors.Is(err, simulated) {
		t.Errorf("expected simulated error on Put, got %v", err)
	}
	if _, err := cache.Exists(ctx, "k"); !errors.Is(err, simulated) {
		t.Errorf("expected simulated error on Exists, got %v", err)
	}
	if err := cache.Delete(ctx, "k"); !errors.Is(err, simulated) {
		t.Errorf("expected simulated error on Delete, got %v", err)
	}
}
