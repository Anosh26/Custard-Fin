package tunnel

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
)

// HTTPTunnelPuller pulls chunks from edge nodes through an HTTP reverse tunnel or local relay.
type HTTPTunnelPuller struct {
	client     *http.Client
	baseURL    string // e.g. "http://localhost:8088" or "http://frps:8000"
	urlPattern string // optional pattern e.g. "http://%s.relay.custardfin.local:8000"
}

// Config defines settings for the HTTP tunnel puller.
type Config struct {
	BaseURL    string        // Base URL of the FRP relay or edge HTTP server
	URLPattern string        // Optional custom URL pattern containing %s for nodeID
	Timeout    time.Duration // HTTP request timeout (default: 10s)
}

// NewHTTPTunnelPuller creates a new HTTPTunnelPuller.
func NewHTTPTunnelPuller(cfg Config) (*HTTPTunnelPuller, error) {
	if strings.TrimSpace(cfg.BaseURL) == "" && strings.TrimSpace(cfg.URLPattern) == "" {
		return nil, errors.New("either BaseURL or URLPattern must be specified")
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &HTTPTunnelPuller{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		urlPattern: cfg.URLPattern,
	}, nil
}

// PullChunk requests a media chunk or playlist from the edge node.
func (p *HTTPTunnelPuller) PullChunk(ctx context.Context, nodeID, mediaID, chunkPath string) (io.ReadCloser, string, int64, error) {
	if strings.TrimSpace(nodeID) == "" {
		return nil, "", 0, errors.New("node id is required")
	}
	if strings.TrimSpace(mediaID) == "" {
		return nil, "", 0, errors.New("media id is required")
	}
	if strings.TrimSpace(chunkPath) == "" {
		return nil, "", 0, errors.New("chunk path is required")
	}

	cleanPath := strings.TrimPrefix(chunkPath, "/")

	var targetURL string
	if p.urlPattern != "" {
		host := fmt.Sprintf(p.urlPattern, nodeID)
		targetURL = fmt.Sprintf("%s/stream/%s/%s", strings.TrimRight(host, "/"), mediaID, cleanPath)
	} else {
		targetURL = fmt.Sprintf("%s/stream/%s/%s", p.baseURL, mediaID, cleanPath)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to create tunnel request: %w", err)
	}

	// For FRP vhost HTTP routing, send the node ID as Host header if using unified baseURL
	if p.urlPattern == "" && p.baseURL != "" {
		req.Header.Set("X-Node-ID", nodeID)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, "", 0, fmt.Errorf("tunnel pull failed for node %s: %w", nodeID, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		_ = resp.Body.Close()
		return nil, "", 0, domain.ErrNotFound
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		_ = resp.Body.Close()
		return nil, "", 0, fmt.Errorf("tunnel returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return resp.Body, contentType, resp.ContentLength, nil
}
