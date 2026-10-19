package tunnel

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
)

func TestHTTPTunnelPuller_Lifecycle(t *testing.T) {
	expectedData := "binary-chunk-data"
	var receivedPath string
	var receivedNodeHeader string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedNodeHeader = r.Header.Get("X-Node-ID")

		if r.URL.Path == "/stream/media1/1080p/seg_0000.ts" {
			w.Header().Set("Content-Type", "video/mp2t")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(expectedData))
			return
		}

		if r.URL.Path == "/stream/media1/missing.ts" {
			http.NotFound(w, r)
			return
		}

		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	puller, err := NewHTTPTunnelPuller(Config{
		BaseURL: ts.URL,
		Timeout: 2 * time.Second,
	})
	if err != nil {
		t.Fatalf("failed to create puller: %v", err)
	}

	ctx := context.Background()

	// 1. Success pull
	body, contentType, length, err := puller.PullChunk(ctx, "node123", "media1", "1080p/seg_0000.ts")
	if err != nil {
		t.Fatalf("failed to pull chunk: %v", err)
	}
	defer body.Close()

	if contentType != "video/mp2t" {
		t.Errorf("expected video/mp2t, got %s", contentType)
	}
	if length != int64(len(expectedData)) {
		t.Errorf("expected length %d, got %d", len(expectedData), length)
	}

	data, err := io.ReadAll(body)
	if err != nil || string(data) != expectedData {
		t.Fatalf("expected data %q, got %q (err: %v)", expectedData, string(data), err)
	}
	if receivedPath != "/stream/media1/1080p/seg_0000.ts" {
		t.Errorf("expected path /stream/media1/1080p/seg_0000.ts, got %s", receivedPath)
	}
	if receivedNodeHeader != "node123" {
		t.Errorf("expected X-Node-ID header node123, got %s", receivedNodeHeader)
	}

	// 2. 404 maps to domain.ErrNotFound
	_, _, _, err = puller.PullChunk(ctx, "node123", "media1", "missing.ts")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for missing chunk, got %v", err)
	}

	// 3. 500 error
	_, _, _, err = puller.PullChunk(ctx, "node123", "media1", "error.ts")
	if err == nil {
		t.Fatalf("expected error on 500 status")
	}
}

func TestHTTPTunnelPuller_Validation(t *testing.T) {
	// Empty config
	_, err := NewHTTPTunnelPuller(Config{})
	if err == nil {
		t.Fatalf("expected error on empty config")
	}

	puller, err := NewHTTPTunnelPuller(Config{BaseURL: "http://localhost:8088"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Empty arguments
	if _, _, _, err := puller.PullChunk(ctx, "", "m", "c"); err == nil {
		t.Errorf("expected error on empty nodeID")
	}
	if _, _, _, err := puller.PullChunk(ctx, "n", "", "c"); err == nil {
		t.Errorf("expected error on empty mediaID")
	}
	if _, _, _, err := puller.PullChunk(ctx, "n", "m", ""); err == nil {
		t.Errorf("expected error on empty chunkPath")
	}
}

func TestHTTPTunnelPuller_ContextCancellation(t *testing.T) {
	// Server that delays response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	puller, err := NewHTTPTunnelPuller(Config{BaseURL: ts.URL, Timeout: 5 * time.Second})
	if err != nil {
		t.Fatalf("failed to create puller: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, _, _, err = puller.PullChunk(ctx, "n", "m", "chunk.ts")
	if err == nil {
		t.Fatalf("expected error on canceled context")
	}
}
