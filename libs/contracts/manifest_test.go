package contracts

import (
	"strings"
	"testing"
	"time"
)

func TestManifestValidation(t *testing.T) {
	valid := &Manifest{
		SchemaVersion: "1.0",
		MediaID:       "550e8400-e29b-41d4-a716-446655440000",
		NodeID:        "550e8400-e29b-41d4-a716-446655440001",
		Title:         "Tears of Steel",
		DurationS:     734,
		ContentHash:   "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Renditions: []Rendition{
			{
				Resolution:     "1080p",
				BitrateBps:     5000000,
				PlaylistPath:   "1080p/index.m3u8",
				ChunkDurationS: 4,
			},
			{
				Resolution:     "720p",
				BitrateBps:     2500000,
				PlaylistPath:   "720p/index.m3u8",
				ChunkDurationS: 4,
			},
		},
		IsPublic:  true,
		UpdatedAt: time.Now().UTC(),
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid manifest, got: %v", err)
	}

	// Test non-4s chunk duration violation
	invalidChunk := *valid
	invalidChunk.Renditions = []Rendition{
		{
			Resolution:     "1080p",
			BitrateBps:     5000000,
			PlaylistPath:   "1080p/index.m3u8",
			ChunkDurationS: 6, // violation
		},
	}
	if err := invalidChunk.Validate(); err == nil {
		t.Fatalf("expected error for non-4s chunk duration")
	}

	// Test size constraint (> 3KB)
	oversized := *valid
	oversized.Description = strings.Repeat("A", 3500)
	if err := oversized.Validate(); err == nil {
		t.Fatalf("expected error for manifest exceeding 3KB")
	}
}

func TestHeartbeatValidation(t *testing.T) {
	hb := &Heartbeat{
		NodeID:        "550e8400-e29b-41d4-a716-446655440001",
		SentAt:        time.Now().UTC(),
		Version:       "0.1.0",
		ActiveStreams: 5,
	}

	if err := hb.Validate(); err != nil {
		t.Fatalf("expected valid heartbeat, got: %v", err)
	}

	// Exceeding per-node cap of 10
	hb.ActiveStreams = 11
	if err := hb.Validate(); err == nil {
		t.Fatalf("expected error for active streams > 10")
	}
}
