package app

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/adapters/fs"
	"github.com/anosh26/custardfin/apps/edge/internal/domain"
	"github.com/anosh26/custardfin/apps/edge/internal/ports"
)

// FakeMediaRepository is an in-memory repository for unit testing.
type FakeMediaRepository struct {
	mu    sync.RWMutex
	items map[string]*domain.Media
}

func NewFakeMediaRepository() *FakeMediaRepository {
	return &FakeMediaRepository{
		items: make(map[string]*domain.Media),
	}
}

func (r *FakeMediaRepository) FindByPath(ctx context.Context, path string) (*domain.Media, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.items[path], nil
}

func (r *FakeMediaRepository) FindByID(ctx context.Context, id string) (*domain.Media, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, m := range r.items {
		if m.ID == id {
			return m, nil
		}
	}
	return nil, nil
}

func (r *FakeMediaRepository) FindByHash(ctx context.Context, hash string) (*domain.Media, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, m := range r.items {
		if m.ContentHash == hash {
			return m, nil
		}
	}
	return nil, nil
}

func (r *FakeMediaRepository) Save(ctx context.Context, media *domain.Media) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[media.Path] = media
	return nil
}

func (r *FakeMediaRepository) List(ctx context.Context) ([]*domain.Media, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.Media
	for _, m := range r.items {
		list = append(list, m)
	}
	return list, nil
}

func (r *FakeMediaRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for p, m := range r.items {
		if m.ID == id {
			delete(r.items, p)
			break
		}
	}
	return nil
}

// FakeProber simulates ffprobe responses.
type FakeProber struct {
	probeCalls int
}

func (p *FakeProber) Probe(ctx context.Context, filePath string) (*ports.ProbeResult, error) {
	p.probeCalls++
	return &ports.ProbeResult{
		DurationS:  120,
		Width:      1920,
		Height:     1080,
		VideoCodec: "h264",
		AudioCodec: "aac",
	}, nil
}

func TestScanService_IncrementalScan(t *testing.T) {
	tempDir := t.TempDir()
	sampleVideo := filepath.Join(tempDir, "big_buck_bunny.mp4")
	if err := os.WriteFile(sampleVideo, []byte("fake mp4 content"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := fs.NewFileScanner()
	prober := &FakeProber{}
	hasher := fs.NewHasher()
	repo := NewFakeMediaRepository()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewScanService(scanner, prober, hasher, repo, logger)
	ctx := context.Background()

	// 1. Initial scan: should detect, hash, probe, and add 1 media
	summary1, err := service.ScanDirectory(ctx, tempDir)
	if err != nil {
		t.Fatalf("first scan failed: %v", err)
	}
	if summary1.Added != 1 || summary1.Unchanged != 0 {
		t.Errorf("expected 1 added and 0 unchanged, got %d added and %d unchanged", summary1.Added, summary1.Unchanged)
	}
	if prober.probeCalls != 1 {
		t.Errorf("expected 1 probe call, got %d", prober.probeCalls)
	}

	// 2. Second scan without modifications: fast-path skip
	summary2, err := service.ScanDirectory(ctx, tempDir)
	if err != nil {
		t.Fatalf("second scan failed: %v", err)
	}
	if summary2.Added != 0 || summary2.Unchanged != 1 {
		t.Errorf("expected 0 added and 1 unchanged, got %d added and %d unchanged", summary2.Added, summary2.Unchanged)
	}
	// probeCalls must still be 1 (skipped ffprobe)
	if prober.probeCalls != 1 {
		t.Errorf("expected probe calls to stay 1 after unchanged scan, got %d", prober.probeCalls)
	}

	// 3. Modify the file: update mod time and size
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(sampleVideo, []byte("fake mp4 content with extra bytes"), 0644); err != nil {
		t.Fatal(err)
	}

	summary3, err := service.ScanDirectory(ctx, tempDir)
	if err != nil {
		t.Fatalf("third scan failed: %v", err)
	}
	if summary3.Updated != 1 {
		t.Errorf("expected 1 updated media, got %d", summary3.Updated)
	}
	if prober.probeCalls != 2 {
		t.Errorf("expected probe calls to increment to 2, got %d", prober.probeCalls)
	}
}
