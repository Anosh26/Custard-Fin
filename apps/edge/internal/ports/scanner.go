package ports

import (
	"context"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/domain"
)

// FileEntry represents a discovered media file on disk before probing.
type FileEntry struct {
	Path    string
	Size    int64
	ModTime time.Time
}

// Scanner defines the interface for walking directories and emitting discovered media files.
type Scanner interface {
	Scan(ctx context.Context, rootDir string) (<-chan FileEntry, <-chan error)
}

// ProbeResult contains metadata extracted from a media container via ffprobe.
type ProbeResult struct {
	DurationS  int
	Width      int
	Height     int
	VideoCodec string
	AudioCodec string
}

// Prober defines the interface for probing media files using ffprobe.
type Prober interface {
	Probe(ctx context.Context, filePath string) (*ProbeResult, error)
}

// MediaRepository defines the persistence contract for discovered media in the local SQLite database.
type MediaRepository interface {
	FindByPath(ctx context.Context, path string) (*domain.Media, error)
	FindByID(ctx context.Context, id string) (*domain.Media, error)
	FindByHash(ctx context.Context, hash string) (*domain.Media, error)
	Save(ctx context.Context, media *domain.Media) error
	List(ctx context.Context) ([]*domain.Media, error)
	Delete(ctx context.Context, id string) error
}
