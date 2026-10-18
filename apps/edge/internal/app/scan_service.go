package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/adapters/fs"
	"github.com/anosh26/custardfin/apps/edge/internal/domain"
	"github.com/anosh26/custardfin/apps/edge/internal/ports"
)

// ScanSummary aggregates metrics from a scan run.
type ScanSummary struct {
	TotalScanned int
	Added        int
	Updated      int
	Unchanged    int
	Failed       int
	Duration     time.Duration
}

// ScanService coordinates media discovery, hashing, probing, and persistence.
type ScanService struct {
	scanner ports.Scanner
	prober  ports.Prober
	hasher  *fs.Hasher
	repo    ports.MediaRepository
	logger  *slog.Logger
}

// NewScanService constructs an instance of ScanService.
func NewScanService(
	scanner ports.Scanner,
	prober ports.Prober,
	hasher *fs.Hasher,
	repo ports.MediaRepository,
	logger *slog.Logger,
) *ScanService {
	if logger == nil {
		logger = slog.Default()
	}
	if hasher == nil {
		hasher = fs.NewHasher()
	}
	return &ScanService{
		scanner: scanner,
		prober:  prober,
		hasher:  hasher,
		repo:    repo,
		logger:  logger,
	}
}

// ScanDirectory walks rootDir, detecting new and modified video files.
func (s *ScanService) ScanDirectory(ctx context.Context, rootDir string) (*ScanSummary, error) {
	startTime := time.Now()
	summary := &ScanSummary{}

	s.logger.Info("beginning library scan", "directory", rootDir)

	entryChan, errChan := s.scanner.Scan(ctx, rootDir)

	for entry := range entryChan {
		summary.TotalScanned++

		// Fast-path change detection: check if already in DB
		existing, err := s.repo.FindByPath(ctx, entry.Path)
		if err != nil {
			s.logger.Error("failed checking existing media path", "path", entry.Path, "err", err)
			summary.Failed++
			continue
		}

		if existing != nil && existing.FileSize == entry.Size && existing.ModTime.Equal(entry.ModTime) {
			// Unchanged file: skip expensive hash and ffprobe
			summary.Unchanged++
			continue
		}

		// Compute content hash (streaming SHA-256)
		hash, err := s.hasher.ComputeSHA256(ctx, entry.Path)
		if err != nil {
			s.logger.Error("failed to compute content hash", "path", entry.Path, "err", err)
			summary.Failed++
			continue
		}

		// Probe media stream properties via ffprobe
		probe, err := s.prober.Probe(ctx, entry.Path)
		if err != nil {
			s.logger.Error("failed to probe media file", "path", entry.Path, "err", err)
			summary.Failed++
			continue
		}

		if existing != nil {
			// Update modified media
			existing.FileSize = entry.Size
			existing.ModTime = entry.ModTime
			existing.ContentHash = hash
			existing.DurationS = probe.DurationS
			existing.Width = probe.Width
			existing.Height = probe.Height
			existing.VideoCodec = probe.VideoCodec
			existing.AudioCodec = probe.AudioCodec
			existing.UpdatedAt = time.Now().UTC()

			if err := s.repo.Save(ctx, existing); err != nil {
				s.logger.Error("failed to update media record", "path", entry.Path, "err", err)
				summary.Failed++
				continue
			}
			summary.Updated++
			s.logger.Info("updated modified media", "title", existing.Title, "path", entry.Path)
		} else {
			// Create newly discovered media
			media, err := domain.NewMedia(
				entry.Path,
				entry.Size,
				entry.ModTime,
				hash,
				probe.DurationS,
				probe.Width,
				probe.Height,
				probe.VideoCodec,
				probe.AudioCodec,
			)
			if err != nil {
				s.logger.Error("failed to construct media domain entity", "path", entry.Path, "err", err)
				summary.Failed++
				continue
			}

			if err := s.repo.Save(ctx, media); err != nil {
				s.logger.Error("failed to save new media record", "path", entry.Path, "err", err)
				summary.Failed++
				continue
			}
			summary.Added++
			s.logger.Info("cataloged new media", "title", media.Title, "duration_s", media.DurationS)
		}
	}

	if err := <-errChan; err != nil && !errors.Is(err, context.Canceled) {
		return summary, fmt.Errorf("scan failed on %s: %w", rootDir, err)
	}

	summary.Duration = time.Since(startTime)
	s.logger.Info("scan completed",
		"total", summary.TotalScanned,
		"added", summary.Added,
		"updated", summary.Updated,
		"unchanged", summary.Unchanged,
		"failed", summary.Failed,
		"duration", summary.Duration,
	)

	return summary, nil
}
