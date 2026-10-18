package sync_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	syncAdapter "github.com/anosh26/custardfin/apps/edge/internal/adapters/sync"
	"github.com/anosh26/custardfin/apps/edge/internal/domain"
	"github.com/anosh26/custardfin/libs/contracts"
)

func validTestMedia(t *testing.T) *domain.Media {
	t.Helper()
	m, err := domain.NewMedia(
		"/media/test_clip.mp4",
		1024000,
		time.Now().UTC(),
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		120,
		1920,
		1080,
		"h264",
		"aac",
	)
	if err != nil {
		t.Fatalf("failed creating test media: %v", err)
	}
	return m
}

func validTestRenditions() []contracts.Rendition {
	return []contracts.Rendition{
		{Resolution: "1080p", BitrateBps: 5000000, PlaylistPath: "1080p/index.m3u8", ChunkDurationS: 4},
		{Resolution: "720p", BitrateBps: 2800000, PlaylistPath: "720p/index.m3u8", ChunkDurationS: 4},
	}
}

func TestNewHTTPManifestSyncer_Validation(t *testing.T) {
	_, err := syncAdapter.NewHTTPManifestSyncer("http://localhost:8080", "", "token", nil)
	if !errors.Is(err, syncAdapter.ErrNodeConfigMissing) {
		t.Errorf("expected ErrNodeConfigMissing for empty nodeID, got %v", err)
	}

	_, err = syncAdapter.NewHTTPManifestSyncer("http://localhost:8080", "node-1", "", nil)
	if !errors.Is(err, syncAdapter.ErrNodeConfigMissing) {
		t.Errorf("expected ErrNodeConfigMissing for empty nodeToken, got %v", err)
	}

	syncer, err := syncAdapter.NewHTTPManifestSyncer("", "node-1", "token-1", nil)
	if err != nil {
		t.Fatalf("unexpected error with default cloud URL: %v", err)
	}
	if syncer == nil {
		t.Fatalf("expected non-nil syncer")
	}
}

func TestHTTPManifestSyncer_SyncManifest_Success(t *testing.T) {
	nodeID := "node-uuid-123"
	nodeToken := "secret-edge-token-xyz"
	media := validTestMedia(t)
	renditions := validTestRenditions()

	var receivedMethod, receivedPath, receivedToken string
	var receivedManifest contracts.Manifest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		receivedToken = r.Header.Get("X-Node-Token")

		_ = json.NewDecoder(r.Body).Decode(&receivedManifest)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"synced"}`))
	}))
	defer server.Close()

	syncer, err := syncAdapter.NewHTTPManifestSyncer(server.URL, nodeID, nodeToken, server.Client())
	if err != nil {
		t.Fatalf("failed creating syncer: %v", err)
	}

	err = syncer.SyncManifest(context.Background(), media, renditions)
	if err != nil {
		t.Fatalf("SyncManifest failed: %v", err)
	}

	if receivedMethod != http.MethodPut {
		t.Errorf("expected PUT method, got %s", receivedMethod)
	}
	expectedPath := "/v1/nodes/" + nodeID + "/media/" + media.ID + "/manifest"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
	if receivedToken != nodeToken {
		t.Errorf("expected token %s, got %s", nodeToken, receivedToken)
	}
	if receivedManifest.Title != media.Title {
		t.Errorf("expected manifest title %s, got %s", media.Title, receivedManifest.Title)
	}
	if len(receivedManifest.Renditions) != 2 {
		t.Errorf("expected 2 renditions, got %d", len(receivedManifest.Renditions))
	}
}

func TestHTTPManifestSyncer_SyncManifest_InvariantRejection(t *testing.T) {
	nodeID := "node-uuid-123"
	nodeToken := "secret-edge-token-xyz"
	media := validTestMedia(t)

	syncer, _ := syncAdapter.NewHTTPManifestSyncer("http://localhost:8080", nodeID, nodeToken, nil)

	// Non-4-second chunk duration violates Invariant 6
	invalidRenditions := []contracts.Rendition{
		{Resolution: "1080p", BitrateBps: 5000000, PlaylistPath: "1080p/index.m3u8", ChunkDurationS: 6},
	}

	err := syncer.SyncManifest(context.Background(), media, invalidRenditions)
	if err == nil {
		t.Fatalf("expected validation error for chunk duration != 4s, got nil")
	}
}

func TestHTTPManifestSyncer_UnsyncManifest_Success(t *testing.T) {
	nodeID := "node-uuid-123"
	nodeToken := "secret-edge-token-xyz"
	mediaID := "media-uuid-456"

	var receivedMethod, receivedPath, receivedToken string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		receivedToken = r.Header.Get("X-Node-Token")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	syncer, _ := syncAdapter.NewHTTPManifestSyncer(server.URL, nodeID, nodeToken, server.Client())

	err := syncer.UnsyncManifest(context.Background(), mediaID)
	if err != nil {
		t.Fatalf("UnsyncManifest failed: %v", err)
	}

	if receivedMethod != http.MethodDelete {
		t.Errorf("expected DELETE method, got %s", receivedMethod)
	}
	expectedPath := "/v1/nodes/" + nodeID + "/media/" + mediaID + "/manifest"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
	if receivedToken != nodeToken {
		t.Errorf("expected token %s, got %s", nodeToken, receivedToken)
	}
}

func TestHTTPManifestSyncer_CloudApiErrorHandling(t *testing.T) {
	nodeID := "node-uuid-123"
	nodeToken := "secret-edge-token-xyz"
	media := validTestMedia(t)
	renditions := validTestRenditions()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"code":"UNAUTHORIZED","message":"Invalid token"}`))
	}))
	defer server.Close()

	syncer, _ := syncAdapter.NewHTTPManifestSyncer(server.URL, nodeID, nodeToken, server.Client())

	err := syncer.SyncManifest(context.Background(), media, renditions)
	if err == nil {
		t.Fatalf("expected error for 401 response")
	}
	if !errors.Is(err, syncAdapter.ErrSyncFailed) {
		t.Errorf("expected error wrapping ErrSyncFailed, got %v", err)
	}
}
