package fs

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestNewHasher(t *testing.T) {
	hasher := NewHasher()
	if hasher == nil {
		t.Fatal("expected non-nil Hasher from NewHasher")
	}
	if hasher.bufferSize != DefaultBufferSize {
		t.Errorf("expected buffer size %d, got %d", DefaultBufferSize, hasher.bufferSize)
	}
}

func TestHasher_ComputeSHA256_StandardFile(t *testing.T) {
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

func TestHasher_ComputeSHA256_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "empty.bin")

	if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
		t.Fatalf("failed to create empty temp file: %v", err)
	}

	// SHA-256 of empty string is e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
	const emptyHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	hasher := NewHasher()
	gotHash, err := hasher.ComputeSHA256(context.Background(), filePath)
	if err != nil {
		t.Fatalf("ComputeSHA256 returned error on empty file: %v", err)
	}

	if gotHash != emptyHash {
		t.Errorf("expected empty hash %q, got %q", emptyHash, gotHash)
	}
}

func TestHasher_ComputeSHA256_MultiBufferFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "multi_block.bin")

	// Create 200 KB file (crosses multiple 64 KB buffer iterations)
	content := make([]byte, 200*1024)
	for i := range content {
		content[i] = byte(i % 251)
	}
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("failed to write multi-block file: %v", err)
	}

	expectedHashRaw := sha256.Sum256(content)
	expectedHash := hex.EncodeToString(expectedHashRaw[:])

	hasher := NewHasher()
	gotHash, err := hasher.ComputeSHA256(context.Background(), filePath)
	if err != nil {
		t.Fatalf("ComputeSHA256 failed on multi-block file: %v", err)
	}

	if gotHash != expectedHash {
		t.Errorf("expected hash %q, got %q", expectedHash, gotHash)
	}
}

func TestHasher_ComputeSHA256_NonExistentFile(t *testing.T) {
	hasher := NewHasher()
	_, err := hasher.ComputeSHA256(context.Background(), "path_that_does_not_exist_at_all.bin")
	if err == nil {
		t.Fatalf("expected error on non-existent file, got nil")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist wrapped, got: %v", err)
	}
}

func TestHasher_ContextCancellation(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "large.bin")

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
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
}
