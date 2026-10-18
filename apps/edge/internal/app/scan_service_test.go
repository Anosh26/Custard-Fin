package app

import (
	"context"
	"errors"
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

// FakeMediaRepository is an in-memory repository for unit testing with failure injection.
type FakeMediaRepository struct {
	mu           sync.RWMutex
	items        map[string]*domain.Media
	failFindByPath bool
	failSave       bool
}

func NewFakeMediaRepository() *FakeMediaRepository {
	return &FakeMediaRepository{
		items: make(map[string]*domain.Media),
	}
}

func (r *FakeMediaRepository) FindByPath(ctx context.Context, path string) (*domain.Media, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.failFindByPath {
		return nil, errors.New("simulated find error")
	}
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
	if r.failSave {
		return errors.New("simulated save error")
	}
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

func (r *FakeMediaRepository) SetPublic(ctx context.Context, id string, isPublic bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, m := range r.items {
		if m.ID == id {
			m.IsPublic = isPublic
			return nil
		}
	}
	return errors.New("media not found")
}

// FakeProber simulates ffprobe responses with failure injection.
type FakeProber struct {
	probeCalls int
	failOn     string
}

func (p *FakeProber) Probe(ctx context.Context, filePath string) (*ports.ProbeResult, error) {
	p.probeCalls++
	if p.failOn != "" && filepath.Base(filePath) == p.failOn {
		return nil, errors.New("simulated probe failure")
	}
	return &ports.ProbeResult{
		DurationS:  120,
		Width:      1920,
		Height:     1080,
		VideoCodec: "h264",
		AudioCodec: "aac",
	}, nil
}

func TestNewScanService_Defaults(t *testing.T) {
	service := NewScanService(fs.NewFileScanner(), &FakeProber{}, nil, NewFakeMediaRepository(), nil)
	if service == nil {
		t.Fatal("expected non-nil ScanService")
	}
	if service.logger == nil {
		t.Error("expected default logger to be initialized")
	}
	if service.hasher == nil {
		t.Error("expected default hasher to be initialized")
	}
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

func TestScanService_FaultTolerance(t *testing.T) {
	tempDir := t.TempDir()

	file1 := filepath.Join(tempDir, "video1.mp4")
	file2 := filepath.Join(tempDir, "failing_probe.mp4")
	file3 := filepath.Join(tempDir, "video3.mp4")

	for _, f := range []string{file1, file2, file3} {
		if err := os.WriteFile(f, []byte("dummy video"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	scanner := fs.NewFileScanner()
	prober := &FakeProber{failOn: "failing_probe.mp4"}
	repo := NewFakeMediaRepository()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewScanService(scanner, prober, nil, repo, logger)
	summary, err := service.ScanDirectory(context.Background(), tempDir)
	if err != nil {
		t.Fatalf("scan directory failed: %v", err)
	}

	// 2 should succeed, 1 failed probe
	if summary.Added != 2 {
		t.Errorf("expected 2 added, got %d", summary.Added)
	}
	if summary.Failed != 1 {
		t.Errorf("expected 1 failure, got %d", summary.Failed)
	}
}

func TestScanService_RepositoryErrors(t *testing.T) {
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "video.mp4")
	_ = os.WriteFile(file, []byte("dummy"), 0644)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// Case 1: FindByPath fails
	repo1 := NewFakeMediaRepository()
	repo1.failFindByPath = true
	service1 := NewScanService(fs.NewFileScanner(), &FakeProber{}, nil, repo1, logger)
	summary1, err := service1.ScanDirectory(context.Background(), tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary1.Failed != 1 {
		t.Errorf("expected 1 failure on FindByPath error, got %d", summary1.Failed)
	}

	// Case 2: Save fails
	repo2 := NewFakeMediaRepository()
	repo2.failSave = true
	service2 := NewScanService(fs.NewFileScanner(), &FakeProber{}, nil, repo2, logger)
	summary2, err := service2.ScanDirectory(context.Background(), tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary2.Failed != 1 {
		t.Errorf("expected 1 failure on Save error, got %d", summary2.Failed)
	}
}
