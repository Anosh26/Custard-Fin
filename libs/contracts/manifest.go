package contracts

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// MaxManifestSizeBytes defines the hard invariant cap on manifest size (3 KB).
const MaxManifestSizeBytes = 3072

// Rendition describes an encoded HLS profile.
type Rendition struct {
	Resolution     string `json:"resolution"`       // e.g. "1080p", "720p", "480p"
	BitrateBps     int64  `json:"bitrate_bps"`      // e.g. 5000000
	PlaylistPath   string `json:"playlist_path"`    // e.g. "1080p/index.m3u8"
	ChunkDurationS int    `json:"chunk_duration_s"` // invariant: 4s segments
}

// Manifest represents the metadata document synced to the cloud catalog.
type Manifest struct {
	SchemaVersion string      `json:"schema_version"`
	MediaID       string      `json:"media_id"`
	NodeID        string      `json:"node_id"`
	Title         string      `json:"title"`
	Description   string      `json:"description,omitempty"`
	DurationS     int         `json:"duration_s"`
	ContentHash   string      `json:"content_hash"` // SHA-256
	Renditions    []Rendition `json:"renditions"`
	IsPublic      bool        `json:"is_public"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

var (
	ErrManifestTooLarge     = errors.New("manifest exceeds maximum size limit of 3 KB")
	ErrInvalidChunkDuration = errors.New("chunk duration must be exactly 4 seconds")
	ErrMissingRequiredField = errors.New("missing required manifest field")
)

// Validate validates the manifest against domain invariants.
func (m *Manifest) Validate() error {
	if m.SchemaVersion != "1.0" {
		return fmt.Errorf("%w: unsupported schema_version %q", ErrMissingRequiredField, m.SchemaVersion)
	}
	if m.MediaID == "" || m.NodeID == "" || m.Title == "" || m.ContentHash == "" {
		return fmt.Errorf("%w: media_id, node_id, title, content_hash are required", ErrMissingRequiredField)
	}
	if m.DurationS <= 0 {
		return fmt.Errorf("%w: duration_s must be positive", ErrMissingRequiredField)
	}
	if len(m.Renditions) == 0 {
		return fmt.Errorf("%w: at least one rendition is required", ErrMissingRequiredField)
	}
	for _, r := range m.Renditions {
		if r.ChunkDurationS != 4 {
			return fmt.Errorf("%w: got %d seconds", ErrInvalidChunkDuration, r.ChunkDurationS)
		}
		if r.PlaylistPath == "" || r.BitrateBps <= 0 {
			return fmt.Errorf("%w: valid playlist_path and positive bitrate_bps required", ErrMissingRequiredField)
		}
	}

	payload, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}
	if len(payload) > MaxManifestSizeBytes {
		return fmt.Errorf("%w: payload is %d bytes", ErrManifestTooLarge, len(payload))
	}

	return nil
}
