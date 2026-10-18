package fs

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestHasher_ComputeSHA256(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "sample.bin")

	content := []byte("CustardFin media content test bytes for streaming hasher")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("failed to create temp test file: %v", err)
	}

	expectedHashRaw := sha256.Sum256(content)
	expectedHash := hex.EncodeToString(expectedHashRaw[:])

	hasher := NewHasher()
	gotHash, err := hasher.ComputeSHA256(context.Background(), filePath)
	if err != nil {
		t.Fatalf("ComputeSHA256 returned error: %v", err)
	}

	if gotHash != expectedHash {
		t.Errorf("expected hash %q, got %q", expectedHash, gotHash)
	}
}

func TestHasher_ContextCancellation(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "large.bin")

	// 1 MB dummy file
	dummy := make([]byte, 1024*1024)
	if err := os.WriteFile(filePath, dummy, 0644); err != nil {
		t.Fatalf("failed to write dummy file: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	hasher := NewHasher()
	_, err := hasher.ComputeSHA256(ctx, filePath)
	if err == nil {
		t.Errorf("expected context cancellation error, got nil")
	}
}
