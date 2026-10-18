package fs

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// DefaultBufferSize is the chunk size used for streaming file hashing (64 KB).
const DefaultBufferSize = 64 * 1024

// Hasher provides streaming content hashing for media files.
type Hasher struct {
	bufferSize int
}

// NewHasher returns a new Hasher instance with default buffer settings.
func NewHasher() *Hasher {
	return &Hasher{
		bufferSize: DefaultBufferSize,
	}
}

// ComputeSHA256 streams the file at filePath in fixed-size chunks and calculates its SHA-256 hash.
// It respects context cancellation for large files.
func (h *Hasher) ComputeSHA256(ctx context.Context, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file for hashing: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	buf := make([]byte, h.bufferSize)

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		n, err := file.Read(buf)
		if n > 0 {
			hasher.Write(buf[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("read error during hashing: %w", err)
		}
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
