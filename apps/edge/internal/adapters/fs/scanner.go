package fs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/anosh26/custardfin/apps/edge/internal/ports"
)

// SupportedExtensions lists the video container formats processed by CustardFin.
var SupportedExtensions = map[string]struct{}{
	".mp4":  {},
	".mkv":  {},
	".mov":  {},
	".webm": {},
	".avi":  {},
}

var (
	ErrInvalidRootDir = errors.New("invalid or empty root directory")
	ErrPathEscaped    = errors.New("discovered path escaped root directory")
)

// FileScanner implements ports.Scanner using standard library filepath walking.
type FileScanner struct {
	allowedExtensions map[string]struct{}
}

// NewFileScanner constructs a FileScanner with supported video extensions.
func NewFileScanner() *FileScanner {
	return &FileScanner{
		allowedExtensions: SupportedExtensions,
	}
}

// Scan traverses rootDir recursively and streams discovered video files on the returned channel.
func (s *FileScanner) Scan(ctx context.Context, rootDir string) (<-chan ports.FileEntry, <-chan error) {
	entryChan := make(chan ports.FileEntry)
	errChan := make(chan error, 1)

	go func() {
		defer close(entryChan)
		defer close(errChan)

		cleanRoot := filepath.Clean(rootDir)
		if cleanRoot == "" || cleanRoot == "." {
			// check if rootDir was really empty
			if strings.TrimSpace(rootDir) == "" {
				errChan <- ErrInvalidRootDir
				return
			}
		}

		walkErr := filepath.WalkDir(cleanRoot, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Check cancellation
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if d.IsDir() {
				return nil
			}

			// Path traversal check
			cleanPath := filepath.Clean(path)
			rel, err := filepath.Rel(cleanRoot, cleanPath)
			if err != nil || strings.HasPrefix(rel, "..") {
				return fmt.Errorf("%w: %s", ErrPathEscaped, cleanPath)
			}

			// Extension check
			ext := strings.ToLower(filepath.Ext(cleanPath))
			if _, ok := s.allowedExtensions[ext]; !ok {
				return nil
			}

			info, err := d.Info()
			if err != nil {
				return fmt.Errorf("failed to get file info for %s: %w", cleanPath, err)
			}

			select {
			case entryChan <- ports.FileEntry{
				Path:    cleanPath,
				Size:    info.Size(),
				ModTime: info.ModTime(),
			}:
			case <-ctx.Done():
				return ctx.Err()
			}

			return nil
		})

		if walkErr != nil && !errors.Is(walkErr, context.Canceled) {
			errChan <- walkErr
		}
	}()

	return entryChan, errChan
}
