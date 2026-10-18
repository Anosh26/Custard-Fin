package heartbeat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/ports"
	"github.com/anosh26/custardfin/libs/contracts"
)

var (
	ErrMissingConfig    = errors.New("node ID and token are required for heartbeat worker")
	ErrHeartbeatFailed  = errors.New("heartbeat dispatch failed")
	ErrWorkerRunning    = errors.New("heartbeat worker is already running")
	ErrWorkerNotRunning = errors.New("heartbeat worker is not running")
)

var _ ports.HeartbeatReporter = (*Worker)(nil)

// Config configures the edge node background heartbeat worker.
type Config struct {
	CloudBaseURL  string
	NodeID        string
	NodeToken     string
	Version       string
	Interval      time.Duration
	GPUCaps       *contracts.GPUCapabilities
	StreamTracker ports.StreamTracker
	HTTPClient    *http.Client
	Logger        *slog.Logger
}

// Worker periodically transmits node liveness, stream load, and hardware capabilities to the Cloud Control API.
type Worker struct {
	cfg        Config
	httpClient *http.Client
	logger     *slog.Logger

	mu      sync.Mutex
	running bool
	stopCh  chan struct{}
	doneCh  chan struct{}
}

// NewWorker constructs a new edge node heartbeat worker.
func NewWorker(cfg Config) (*Worker, error) {
	cfg.NodeID = strings.TrimSpace(cfg.NodeID)
	cfg.NodeToken = strings.TrimSpace(cfg.NodeToken)
	if cfg.NodeID == "" || cfg.NodeToken == "" {
		return nil, ErrMissingConfig
	}

	cfg.CloudBaseURL = strings.TrimRight(strings.TrimSpace(cfg.CloudBaseURL), "/")
	if cfg.CloudBaseURL == "" {
		cfg.CloudBaseURL = "http://localhost:8080"
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 15 * time.Second
	}
	if cfg.Version == "" {
		cfg.Version = "1.0.0"
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 10 * time.Second}
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return &Worker{
		cfg:        cfg,
		httpClient: cfg.HTTPClient,
		logger:     cfg.Logger,
	}, nil
}

// SendHeartbeat dispatches a single heartbeat payload to the Cloud Control API.
// If hb is nil, an automatic heartbeat is synthesized from current node state and active stream count.
func (w *Worker) SendHeartbeat(ctx context.Context, hb *contracts.Heartbeat) error {
	if hb == nil {
		activeStreams := 0
		if w.cfg.StreamTracker != nil {
			activeStreams = w.cfg.StreamTracker.ActiveStreams()
		}
		if activeStreams < 0 {
			activeStreams = 0
		}
		// Clamp to max allowed invariant cap so heartbeat passes contracts validation
		if activeStreams > contracts.MaxConcurrentStreamsPerNode {
			activeStreams = contracts.MaxConcurrentStreamsPerNode
		}

		hb = &contracts.Heartbeat{
			NodeID:        w.cfg.NodeID,
			SentAt:        time.Now().UTC(),
			Version:       w.cfg.Version,
			GPUCaps:       w.cfg.GPUCaps,
			ActiveStreams: activeStreams,
		}
	}

	if err := hb.Validate(); err != nil {
		return fmt.Errorf("invalid heartbeat payload: %w", err)
	}

	bodyBytes, err := json.Marshal(hb)
	if err != nil {
		return fmt.Errorf("failed marshaling heartbeat JSON: %w", err)
	}

	endpoint := fmt.Sprintf("%s/v1/nodes/%s/heartbeat", w.cfg.CloudBaseURL, w.cfg.NodeID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed creating heartbeat request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Node-Token", w.cfg.NodeToken)

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: outbound HTTP request failed: %v", ErrHeartbeatFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("%w: unexpected response status %d (%s)", ErrHeartbeatFailed, resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	return nil
}

// Start launches the background heartbeat routine.
func (w *Worker) Start(ctx context.Context) error {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return ErrWorkerRunning
	}
	w.running = true
	w.stopCh = make(chan struct{})
	w.doneCh = make(chan struct{})
	w.mu.Unlock()

	w.logger.Info("starting edge heartbeat worker",
		"node_id", w.cfg.NodeID,
		"interval", w.cfg.Interval,
		"cloud_url", w.cfg.CloudBaseURL,
	)

	// Send an initial heartbeat immediately upon starting
	initCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	if err := w.SendHeartbeat(initCtx, nil); err != nil {
		w.logger.Warn("initial heartbeat dispatch failed (cloud may still be starting)", "err", err)
	} else {
		w.logger.Info("initial heartbeat dispatched successfully")
	}
	cancel()

	go w.loop(ctx)
	return nil
}

func (w *Worker) loop(ctx context.Context) {
	defer close(w.doneCh)

	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("context cancelled, terminating heartbeat loop")
			return
		case <-w.stopCh:
			w.logger.Info("stop signal received, terminating heartbeat loop")
			return
		case <-ticker.C:
			hbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			if err := w.SendHeartbeat(hbCtx, nil); err != nil {
				w.logger.Warn("periodic heartbeat dispatch failed", "err", err)
			} else {
				w.logger.Debug("heartbeat dispatched successfully", "sent_at", time.Now().UTC())
			}
			cancel()
		}
	}
}

// Stop terminates the background heartbeat worker cleanly.
func (w *Worker) Stop() error {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return ErrWorkerNotRunning
	}
	w.running = false
	close(w.stopCh)
	w.mu.Unlock()

	select {
	case <-w.doneCh:
		w.logger.Info("edge heartbeat worker stopped successfully")
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("heartbeat worker shutdown timed out")
	}
}
