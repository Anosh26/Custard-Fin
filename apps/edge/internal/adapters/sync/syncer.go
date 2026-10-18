package sync

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/domain"
	"github.com/anosh26/custardfin/apps/edge/internal/ports"
	"github.com/anosh26/custardfin/libs/contracts"
)

var (
	ErrNodeConfigMissing = errors.New("node ID and token are required for manifest sync")
	ErrSyncFailed        = errors.New("cloud sync failed")
)

type HTTPManifestSyncer struct {
	cloudBaseURL string
	nodeID       string
	nodeToken    string
	client       *http.Client
}

// NewHTTPManifestSyncer constructs a new HTTPManifestSyncer adapter.
func NewHTTPManifestSyncer(cloudBaseURL, nodeID, nodeToken string, client *http.Client) (ports.ManifestSyncer, error) {
	cloudBaseURL = strings.TrimRight(strings.TrimSpace(cloudBaseURL), "/")
	if cloudBaseURL == "" {
		cloudBaseURL = "http://localhost:8080"
	}
	if strings.TrimSpace(nodeID) == "" || strings.TrimSpace(nodeToken) == "" {
		return nil, ErrNodeConfigMissing
	}
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	return &HTTPManifestSyncer{
		cloudBaseURL: cloudBaseURL,
		nodeID:       nodeID,
		nodeToken:    nodeToken,
		client:       client,
	}, nil
}

func (s *HTTPManifestSyncer) SyncManifest(ctx context.Context, media *domain.Media, renditions []contracts.Rendition) error {
	if media == nil {
		return fmt.Errorf("%w: media is nil", domain.ErrEmptyPath)
	}

	manifest := &contracts.Manifest{
		SchemaVersion: "1.0",
		MediaID:       media.ID,
		NodeID:        s.nodeID,
		Title:         media.Title,
		DurationS:     media.DurationS,
		ContentHash:   media.ContentHash,
		Renditions:    renditions,
		IsPublic:      true,
		UpdatedAt:     time.Now().UTC(),
	}

	// Validate against domain invariants (chunk duration == 4s, size <= 3 KB)
	if err := manifest.Validate(); err != nil {
		return fmt.Errorf("manifest validation failed: %w", err)
	}

	bodyBytes, err := json.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest JSON: %w", err)
	}

	// Enforce strict 3 KB limit before making outbound call
	if len(bodyBytes) > contracts.MaxManifestSizeBytes {
		return fmt.Errorf("%w: manifest size %d exceeds 3 KB cap", contracts.ErrManifestTooLarge, len(bodyBytes))
	}

	endpoint := fmt.Sprintf("%s/v1/nodes/%s/media/%s/manifest", s.cloudBaseURL, s.nodeID, media.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed creating sync request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Node-Token", s.nodeToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: outbound request failed: %v", ErrSyncFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("%w: cloud api returned status %d: %s", ErrSyncFailed, resp.StatusCode, string(respBody))
	}

	return nil
}

func (s *HTTPManifestSyncer) UnsyncManifest(ctx context.Context, mediaID string) error {
	mediaID = strings.TrimSpace(mediaID)
	if mediaID == "" {
		return errors.New("media ID is required for unsync")
	}

	endpoint := fmt.Sprintf("%s/v1/nodes/%s/media/%s/manifest", s.cloudBaseURL, s.nodeID, mediaID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed creating unsync request: %w", err)
	}

	req.Header.Set("X-Node-Token", s.nodeToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: outbound delete failed: %v", ErrSyncFailed, err)
	}
	defer resp.Body.Close()

	// 204 No Content is success; 404 is also acceptable (already absent)
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("%w: cloud api returned status %d: %s", ErrSyncFailed, resp.StatusCode, string(respBody))
	}

	return nil
}
