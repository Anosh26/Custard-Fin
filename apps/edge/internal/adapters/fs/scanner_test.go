package fs

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileScanner(t *testing.T) {
	scanner := NewFileScanner()
	if scanner == nil {
		t.Fatal("expected non-nil FileScanner from NewFileScanner")
	}
	if len(scanner.allowedExtensions) != len(SupportedExtensions) {
		t.Errorf("expected %d extensions, got %d", len(SupportedExtensions), len(scanner.allowedExtensions))
	}
}

func TestFileScanner_Scan_AllSupportedExtensionsAndCasing(t *testing.T) {
	tempDir := t.TempDir()

	filesToCreate := []string{
		"video1.mp4",
		"video2.MP4",
		"video3.mkv",
		"video4.MKV",
		"video5.mov",
		"video6.MOV",
		"video7.webm",
		"video8.WEBM",
		"video9.avi",
		"video10.AVI",
		// Ignored non-video files:
		"subtitles.srt",
		"cover.jpg",
		"info.nfo",
		"notes.txt",
		"script.sh",
	}

	for _, fname := range filesToCreate {
		if err := os.WriteFile(filepath.Join(tempDir, fname), []byte("dummy data"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	scanner := NewFileScanner()
	ctx := context.Background()
	entryChan, errChan := scanner.Scan(ctx, tempDir)

	var discovered []string
	for entry := range entryChan {
		discovered = append(discovered, filepath.Base(entry.Path))
	}

	if err := <-errChan; err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}

	if len(discovered) != 10 {
		t.Fatalf("expected 10 video files discovered, got %d: %v", len(discovered), discovered)
	}
}

func TestFileScanner_Scan_NestedSubdirectories(t *testing.T) {
	tempDir := t.TempDir()

	// Tree: tempDir/season1/disc1/episode.mkv
	nestedDir := filepath.Join(tempDir, "season1", "disc1")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatal(err)
	}
	targetFile := filepath.Join(nestedDir, "episode.mkv")
	if err := os.WriteFile(targetFile, []byte("episode content"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewFileScanner()
	entryChan, errChan := scanner.Scan(context.Background(), tempDir)

	var discovered []string
	for entry := range entryChan {
		discovered = append(discovered, entry.Path)
	}

	if err := <-errChan; err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}

	if len(discovered) != 1 {
		t.Fatalf("expected 1 nested file, got %d", len(discovered))
	}
	if filepath.Clean(discovered[0]) != filepath.Clean(targetFile) {
		t.Errorf("expected %s, got %s", targetFile, discovered[0])
	}
}

func TestFileScanner_Scan_EmptyAndWhitespaceRootDir(t *testing.T) {
	scanner := NewFileScanner()

	for _, emptyPath := range []string{"", "   ", "\t"} {
		_, errChan := scanner.Scan(context.Background(), emptyPath)
		err := <-errChan
		if !errors.Is(err, ErrInvalidRootDir) {
			t.Errorf("expected ErrInvalidRootDir for root %q, got: %v", emptyPath, err)
		}
	}
}

func TestFileScanner_Scan_NonExistentRootDir(t *testing.T) {
	scanner := NewFileScanner()
	_, errChan := scanner.Scan(context.Background(), "path_that_definitely_does_not_exist_12345")
	err := <-errChan
	if err == nil {
		t.Fatal("expected error on non-existent root dir, got nil")
	}
}

func TestFileScanner_Scan_ContextCancellation(t *testing.T) {
	tempDir := t.TempDir()
	for i := 0; i < 20; i++ {
		fname := filepath.Join(tempDir, "video.mp4")
		_ = os.WriteFile(fname, []byte("x"), 0644)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	scanner := NewFileScanner()
	entryChan, errChan := scanner.Scan(ctx, tempDir)

	for range entryChan {
	}

	err := <-errChan
	// Cancelled context should either return nil or context.Canceled without hanging
	if err != nil && !errors.Is(err, context.Canceled) {
		t.Errorf("unexpected error on cancellation: %v", err)
	}
}
