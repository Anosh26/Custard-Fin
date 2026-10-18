package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/anosh26/custardfin/libs/contracts"
)

// Rendition represents a transcoded quality tier of a media item in the catalog.
type Rendition struct {
	Resolution     string `json:"resolution"`
	BitrateBps     int    `json:"bitrate_bps"`
	PlaylistPath   string `json:"playlist_path"`
	ChunkDurationS int    `json:"chunk_duration_s"`
}

// CatalogItem represents a published media item catalog entry in the cloud database.
// Video binaries remain strictly on the edge node; only this metadata and manifest are stored.
type CatalogItem struct {
	MediaID      string      `json:"media_id"`
	NodeID       string      `json:"node_id"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	DurationS    int         `json:"duration_s"`
	ContentHash  string      `json:"content_hash"`
	IsPublic     bool        `json:"is_public"`
	ManifestJSON []byte      `json:"-"`
	Renditions   []Rendition `json:"renditions"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// NewCatalogItem constructs and validates a new CatalogItem entity.
func NewCatalogItem(
	mediaID, nodeID, title, description string,
	durationS int,
	contentHash string,
	isPublic bool,
	manifestJSON []byte,
	renditions []Rendition,
) (*CatalogItem, error) {
	mediaID = strings.TrimSpace(mediaID)
	if mediaID == "" {
		return nil, fmt.Errorf("%w: media ID cannot be empty", ErrValidation)
	}

	nodeID = strings.TrimSpace(nodeID)
	if nodeID == "" {
		return nil, fmt.Errorf("%w: node ID cannot be empty", ErrValidation)
	}

	title = strings.TrimSpace(title)
	if title == "" || len(title) > 255 {
		return nil, fmt.Errorf("%w: title must be between 1 and 255 characters", ErrValidation)
	}

	if durationS <= 0 {
		return nil, fmt.Errorf("%w: duration must be greater than 0", ErrValidation)
	}

	contentHash = strings.TrimSpace(contentHash)
	if contentHash == "" {
		return nil, fmt.Errorf("%w: content hash cannot be empty", ErrValidation)
	}

	if len(manifestJSON) == 0 {
		return nil, fmt.Errorf("%w: manifest JSON cannot be empty", ErrValidation)
	}

	if len(manifestJSON) > contracts.MaxManifestSizeBytes {
		return nil, fmt.Errorf("%w: manifest size %d bytes exceeds cap of %d bytes",
			ErrManifestTooLarge, len(manifestJSON), contracts.MaxManifestSizeBytes)
	}

	if len(renditions) == 0 {
		return nil, fmt.Errorf("%w: at least one rendition is required", ErrValidation)
	}

	for i, r := range renditions {
		res := strings.ToLower(strings.TrimSpace(r.Resolution))
		if res != "1080p" && res != "720p" && res != "480p" {
			return nil, fmt.Errorf("%w: rendition [%d] invalid resolution %s", ErrValidation, i, r.Resolution)
		}
		if r.BitrateBps <= 0 {
			return nil, fmt.Errorf("%w: rendition [%d] bitrate must be > 0", ErrValidation, i)
		}
		if strings.TrimSpace(r.PlaylistPath) == "" {
			return nil, fmt.Errorf("%w: rendition [%d] playlist path cannot be empty", ErrValidation, i)
		}
		if r.ChunkDurationS != 4 {
			return nil, fmt.Errorf("%w: rendition [%d] chunk duration must be exactly 4 seconds (got %d)",
				ErrValidation, i, r.ChunkDurationS)
		}
	}

	now := time.Now().UTC()
	return &CatalogItem{
		MediaID:      mediaID,
		NodeID:       nodeID,
		Title:        title,
		Description:  strings.TrimSpace(description),
		DurationS:    durationS,
		ContentHash:  contentHash,
		IsPublic:     isPublic,
		ManifestJSON: manifestJSON,
		Renditions:   renditions,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}
