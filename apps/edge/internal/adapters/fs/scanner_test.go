package fs

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileScanner_Scan(t *testing.T) {
	tempDir := t.TempDir()

	// Setup directory tree:
	// tempDir/
	//   ├── movie1.mp4
	//   ├── ignore.txt
	//   ├── season1/
	//   │   ├── episode1.mkv
	//   │   └── notes.nfo
	//   └── season2/
	//       └── episode2.webm

	if err := os.WriteFile(filepath.Join(tempDir, "movie1.mp4"), []byte("video1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "ignore.txt"), []byte("text"), 0644); err != nil {
		t.Fatal(err)
	}

	s1 := filepath.Join(tempDir, "season1")
	if err := os.MkdirAll(s1, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(s1, "episode1.mkv"), []byte("video2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(s1, "notes.nfo"), []byte("info"), 0644); err != nil {
		t.Fatal(err)
	}

	s2 := filepath.Join(tempDir, "season2")
	if err := os.MkdirAll(s2, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(s2, "episode2.webm"), []byte("video3"), 0644); err != nil {
		t.Fatal(err)
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

	if len(discovered) != 3 {
		t.Fatalf("expected 3 video files, discovered %d: %v", len(discovered), discovered)
	}

	expectedMap := map[string]bool{
		"movie1.mp4":    true,
		"episode1.mkv":  true,
		"episode2.webm": true,
	}

	for _, name := range discovered {
		if !expectedMap[name] {
			t.Errorf("unexpected file discovered: %s", name)
		}
	}
}

func TestFileScanner_EmptyRootDir(t *testing.T) {
	scanner := NewFileScanner()
	_, errChan := scanner.Scan(context.Background(), "")

	err := <-errChan
	if err != ErrInvalidRootDir {
		t.Fatalf("expected ErrInvalidRootDir, got: %v", err)
	}
}
