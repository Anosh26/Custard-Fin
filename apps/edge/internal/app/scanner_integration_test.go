package app

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/adapters/ffprobe"
	"github.com/anosh26/custardfin/apps/edge/internal/adapters/fs"
)

// TestScanner_EndToEndIntegration performs full integration of Scanner, Hasher, Subprocess Prober, and Repository
// with genuine files generated using ffmpeg.
func TestScanner_EndToEndIntegration(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found, skipping end-to-end integration test")
	}
	if _, err := exec.LookPath("ffprobe"); err != nil {
		t.Skip("ffprobe not found, skipping end-to-end integration test")
	}

	tempDir := t.TempDir()

	// 1. Create a directory structure with real videos and non-videos
	moviesDir := filepath.Join(tempDir, "movies")
	seriesDir := filepath.Join(tempDir, "series", "season1")
	if err := os.MkdirAll(moviesDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(seriesDir, 0755); err != nil {
		t.Fatal(err)
	}

	movie1 := filepath.Join(moviesDir, "tears_of_steel.mp4")
	episode1 := filepath.Join(seriesDir, "episode_01.mkv")
	ignoreTxt := filepath.Join(moviesDir, "poster.txt")

	_ = os.WriteFile(ignoreTxt, []byte("ignore me"), 0644)

	// Generate movie1: 2s 640x360 H264 + AAC
	cmd1 := exec.Command("ffmpeg", "-y",
		"-f", "lavfi", "-i", "testsrc=duration=2:size=640x360:rate=30",
		"-f", "lavfi", "-i", "sine=frequency=1000:duration=2",
		"-c:v", "libx264", "-c:a", "aac",
		movie1,
	)
	if out, err := cmd1.CombinedOutput(); err != nil {
		t.Fatalf("failed generating movie1: %v (%s)", err, string(out))
	}

	// Generate episode1: 3s 320x240 H264 + AAC
	cmd2 := exec.Command("ffmpeg", "-y",
		"-f", "lavfi", "-i", "testsrc=duration=3:size=320x240:rate=30",
		"-f", "lavfi", "-i", "sine=frequency=440:duration=3",
		"-c:v", "libx264", "-c:a", "aac",
		episode1,
	)
	if out, err := cmd2.CombinedOutput(); err != nil {
		t.Fatalf("failed generating episode1: %v (%s)", err, string(out))
	}

	// 2. Wire up concrete adapters
	scanner := fs.NewFileScanner()
	hasher := fs.NewHasher()
	prober := ffprobe.NewSubprocessProber("")
	repo := NewFakeMediaRepository()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewScanService(scanner, prober, hasher, repo, logger)
	ctx := context.Background()

	// 3. First scan: should discover and process 2 videos
	summary1, err := service.ScanDirectory(ctx, tempDir)
	if err != nil {
		t.Fatalf("first scan failed: %v", err)
	}

	if summary1.TotalScanned != 2 {
		t.Errorf("expected 2 total scanned files, got %d", summary1.TotalScanned)
	}
	if summary1.Added != 2 {
		t.Errorf("expected 2 added media entities, got %d", summary1.Added)
	}
	if summary1.Failed != 0 {
		t.Errorf("expected 0 failures, got %d", summary1.Failed)
	}

	// Verify cataloged details in repository
	m1, err := repo.FindByPath(ctx, movie1)
	if err != nil || m1 == nil {
		t.Fatalf("expected movie1 in repo, got err: %v", err)
	}
	if m1.Title != "tears_of_steel" {
		t.Errorf("expected title 'tears_of_steel', got %q", m1.Title)
	}
	if m1.DurationS != 2 || m1.Width != 640 || m1.Height != 360 {
		t.Errorf("expected 2s 640x360, got %ds %dx%d", m1.DurationS, m1.Width, m1.Height)
	}
	if len(m1.ContentHash) != 64 {
		t.Errorf("expected 64-char SHA-256 hash, got %q", m1.ContentHash)
	}

	m2, err := repo.FindByPath(ctx, episode1)
	if err != nil || m2 == nil {
		t.Fatalf("expected episode1 in repo, got err: %v", err)
	}
	if m2.Title != "episode_01" {
		t.Errorf("expected title 'episode_01', got %q", m2.Title)
	}
	if m2.DurationS != 3 || m2.Width != 320 || m2.Height != 240 {
		t.Errorf("expected 3s 320x240, got %ds %dx%d", m2.DurationS, m2.Width, m2.Height)
	}

	// 4. Second scan: should immediately skip unchanged files
	summary2, err := service.ScanDirectory(ctx, tempDir)
	if err != nil {
		t.Fatalf("second scan failed: %v", err)
	}
	if summary2.Added != 0 || summary2.Unchanged != 2 {
		t.Errorf("expected 0 added and 2 unchanged on fast skip, got %d added and %d unchanged", summary2.Added, summary2.Unchanged)
	}

	// 5. Append/modify episode1 to trigger update
	time.Sleep(15 * time.Millisecond)
	f, err := os.OpenFile(episode1, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.Write([]byte("extra byte"))
	_ = f.Close()

	summary3, err := service.ScanDirectory(ctx, tempDir)
	if err != nil {
		t.Fatalf("third scan failed: %v", err)
	}
	if summary3.Updated != 1 || summary3.Unchanged != 1 {
		t.Errorf("expected 1 updated and 1 unchanged, got %d updated and %d unchanged", summary3.Updated, summary3.Unchanged)
	}
}
