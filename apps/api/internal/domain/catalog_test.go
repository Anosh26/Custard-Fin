package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/libs/contracts"
)

func validRenditions() []domain.Rendition {
	return []domain.Rendition{
		{
			Resolution:     "1080p",
			BitrateBps:     5000000,
			PlaylistPath:   "1080p.m3u8",
			ChunkDurationS: 4,
		},
		{
			Resolution:     "720p",
			BitrateBps:     2800000,
			PlaylistPath:   "720p.m3u8",
			ChunkDurationS: 4,
		},
	}
}

func TestNewCatalogItem_Success(t *testing.T) {
	manifestBytes := []byte(`{"schema_version":"1.0","title":"Big Buck Bunny"}`)
	item, err := domain.NewCatalogItem(
		"media-123",
		"node-456",
		"Big Buck Bunny",
		"Open source movie",
		600,
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		true,
		manifestBytes,
		validRenditions(),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if item.MediaID != "media-123" {
		t.Errorf("expected mediaID media-123, got %s", item.MediaID)
	}
	if item.NodeID != "node-456" {
		t.Errorf("expected nodeID node-456, got %s", item.NodeID)
	}
	if item.Title != "Big Buck Bunny" {
		t.Errorf("expected title Big Buck Bunny, got %s", item.Title)
	}
	if item.DurationS != 600 {
		t.Errorf("expected duration 600, got %d", item.DurationS)
	}
	if !item.IsPublic {
		t.Errorf("expected isPublic=true")
	}
	if len(item.Renditions) != 2 {
		t.Errorf("expected 2 renditions, got %d", len(item.Renditions))
	}
}

func TestNewCatalogItem_ValidationErrors(t *testing.T) {
	manifestBytes := []byte(`{"schema_version":"1.0"}`)

	tests := []struct {
		name         string
		mediaID      string
		nodeID       string
		title        string
		durationS    int
		contentHash  string
		manifestJSON []byte
		renditions   []domain.Rendition
		wantErr      error
	}{
		{
			name:         "empty mediaID",
			mediaID:      "",
			nodeID:       "node-1",
			title:        "Title",
			durationS:    10,
			contentHash:  "hash",
			manifestJSON: manifestBytes,
			renditions:   validRenditions(),
			wantErr:      domain.ErrValidation,
		},
		{
			name:         "empty nodeID",
			mediaID:      "media-1",
			nodeID:       "",
			title:        "Title",
			durationS:    10,
			contentHash:  "hash",
			manifestJSON: manifestBytes,
			renditions:   validRenditions(),
			wantErr:      domain.ErrValidation,
		},
		{
			name:         "empty title",
			mediaID:      "media-1",
			nodeID:       "node-1",
			title:        "",
			durationS:    10,
			contentHash:  "hash",
			manifestJSON: manifestBytes,
			renditions:   validRenditions(),
			wantErr:      domain.ErrValidation,
		},
		{
			name:         "title too long",
			mediaID:      "media-1",
			nodeID:       "node-1",
			title:        strings.Repeat("a", 256),
			durationS:    10,
			contentHash:  "hash",
			manifestJSON: manifestBytes,
			renditions:   validRenditions(),
			wantErr:      domain.ErrValidation,
		},
		{
			name:         "zero duration",
			mediaID:      "media-1",
			nodeID:       "node-1",
			title:        "Title",
			durationS:    0,
			contentHash:  "hash",
			manifestJSON: manifestBytes,
			renditions:   validRenditions(),
			wantErr:      domain.ErrValidation,
		},
		{
			name:         "negative duration",
			mediaID:      "media-1",
			nodeID:       "node-1",
			title:        "Title",
			durationS:    -5,
			contentHash:  "hash",
			manifestJSON: manifestBytes,
			renditions:   validRenditions(),
			wantErr:      domain.ErrValidation,
		},
		{
			name:         "empty contentHash",
			mediaID:      "media-1",
			nodeID:       "node-1",
			title:        "Title",
			durationS:    10,
			contentHash:  "",
			manifestJSON: manifestBytes,
			renditions:   validRenditions(),
			wantErr:      domain.ErrValidation,
		},
		{
			name:         "empty manifest JSON",
			mediaID:      "media-1",
			nodeID:       "node-1",
			title:        "Title",
			durationS:    10,
			contentHash:  "hash",
			manifestJSON: nil,
			renditions:   validRenditions(),
			wantErr:      domain.ErrValidation,
		},
		{
			name:         "manifest JSON exceeds 3 KB cap",
			mediaID:      "media-1",
			nodeID:       "node-1",
			title:        "Title",
			durationS:    10,
			contentHash:  "hash",
			manifestJSON: make([]byte, contracts.MaxManifestSizeBytes+1),
			renditions:   validRenditions(),
			wantErr:      domain.ErrManifestTooLarge,
		},
		{
			name:         "no renditions",
			mediaID:      "media-1",
			nodeID:       "node-1",
			title:        "Title",
			durationS:    10,
			contentHash:  "hash",
			manifestJSON: manifestBytes,
			renditions:   nil,
			wantErr:      domain.ErrValidation,
		},
		{
			name:        "rendition invalid resolution",
			mediaID:     "media-1",
			nodeID:      "node-1",
			title:       "Title",
			durationS:   10,
			contentHash: "hash",
			manifestJSON: manifestBytes,
			renditions: []domain.Rendition{
				{Resolution: "4k", BitrateBps: 1000000, PlaylistPath: "4k.m3u8", ChunkDurationS: 4},
			},
			wantErr: domain.ErrValidation,
		},
		{
			name:        "rendition zero bitrate",
			mediaID:     "media-1",
			nodeID:      "node-1",
			title:       "Title",
			durationS:   10,
			contentHash: "hash",
			manifestJSON: manifestBytes,
			renditions: []domain.Rendition{
				{Resolution: "1080p", BitrateBps: 0, PlaylistPath: "1080p.m3u8", ChunkDurationS: 4},
			},
			wantErr: domain.ErrValidation,
		},
		{
			name:        "rendition empty playlist path",
			mediaID:     "media-1",
			nodeID:      "node-1",
			title:       "Title",
			durationS:   10,
			contentHash: "hash",
			manifestJSON: manifestBytes,
			renditions: []domain.Rendition{
				{Resolution: "1080p", BitrateBps: 1000000, PlaylistPath: "", ChunkDurationS: 4},
			},
			wantErr: domain.ErrValidation,
		},
		{
			name:        "rendition invalid chunk duration (not 4s)",
			mediaID:     "media-1",
			nodeID:      "node-1",
			title:       "Title",
			durationS:   10,
			contentHash: "hash",
			manifestJSON: manifestBytes,
			renditions: []domain.Rendition{
				{Resolution: "1080p", BitrateBps: 1000000, PlaylistPath: "1080p.m3u8", ChunkDurationS: 6},
			},
			wantErr: domain.ErrValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := domain.NewCatalogItem(
				tt.mediaID,
				tt.nodeID,
				tt.title,
				"description",
				tt.durationS,
				tt.contentHash,
				true,
				tt.manifestJSON,
				tt.renditions,
			)
			if err == nil {
				t.Fatalf("expected error, got item: %+v", item)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error wrapping %v, got %v", tt.wantErr, err)
			}
		})
	}
}
