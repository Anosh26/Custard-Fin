package contracts

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func validTestManifest() *Manifest {
	return &Manifest{
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
}

func TestManifestValidation_HappyPath(t *testing.T) {
	manifest := validTestManifest()
	if err := manifest.Validate(); err != nil {
		t.Fatalf("expected valid manifest to pass, got: %v", err)
	}
}

func TestManifestValidation_SchemaVersion(t *testing.T) {
	manifest := validTestManifest()
	manifest.SchemaVersion = "2.0"
	if err := manifest.Validate(); !errors.Is(err, ErrMissingRequiredField) {
		t.Errorf("expected ErrMissingRequiredField for invalid schema version, got: %v", err)
	}
}

func TestManifestValidation_RequiredFields(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(m *Manifest)
	}{
		{"empty MediaID", func(m *Manifest) { m.MediaID = "" }},
		{"empty NodeID", func(m *Manifest) { m.NodeID = "" }},
		{"empty Title", func(m *Manifest) { m.Title = "" }},
		{"empty ContentHash", func(m *Manifest) { m.ContentHash = "" }},
		{"zero duration", func(m *Manifest) { m.DurationS = 0 }},
		{"negative duration", func(m *Manifest) { m.DurationS = -10 }},
		{"empty renditions", func(m *Manifest) { m.Renditions = nil }},
		{"empty playlist path", func(m *Manifest) { m.Renditions[0].PlaylistPath = "" }},
		{"zero bitrate", func(m *Manifest) { m.Renditions[0].BitrateBps = 0 }},
		{"negative bitrate", func(m *Manifest) { m.Renditions[0].BitrateBps = -500 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := validTestManifest()
			tt.mutate(m)
			if err := m.Validate(); err == nil {
				t.Errorf("expected validation error for %s, got nil", tt.name)
			}
		})
	}
}

func TestManifestValidation_ChunkDurationInvariant(t *testing.T) {
	manifest := validTestManifest()
	manifest.Renditions[0].ChunkDurationS = 6 // must be exactly 4s
	if err := manifest.Validate(); !errors.Is(err, ErrInvalidChunkDuration) {
		t.Errorf("expected ErrInvalidChunkDuration, got: %v", err)
	}
}

func TestManifestValidation_SizeCap(t *testing.T) {
	manifest := validTestManifest()
	manifest.Description = strings.Repeat("X", 3500)
	if err := manifest.Validate(); !errors.Is(err, ErrManifestTooLarge) {
		t.Errorf("expected ErrManifestTooLarge, got: %v", err)
	}
}
